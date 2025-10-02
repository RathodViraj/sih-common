package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/RathodViraj/sih-common/configs"
	"github.com/RathodViraj/sih-common/models"
	"github.com/RathodViraj/sih-common/pkg/firebase"

	"firebase.google.com/go/v4/messaging"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const tokenCollectionName = "fcm_tokens"

type NotificationService struct {
	DB          *mongo.Database
	Redis       *redis.Client
	FCMClient   *messaging.Client
	Concurrency int
}

func NewNotificationService(client *mongo.Client, redisClient *redis.Client) *NotificationService {
	cfg := configs.LoadDBConfig()
	db := client.Database(cfg.Database)

	ctx := context.Background()

	app := firebase.App

	mClient, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting messaging client: %v", err)
	}

	svc := &NotificationService{
		DB:          db,
		Redis:       redisClient,
		FCMClient:   mClient,
		Concurrency: 50,
	}

	svc.StartRetryWorker(ctx)

	return svc
}

func (s *NotificationService) RegisterOrUpdateToken(ctx context.Context, userId, token, platform string) error {
	collection := s.DB.Collection(tokenCollectionName)
	filter := bson.M{"userId": userId, "token": token}
	update := bson.M{
		"$set": bson.M{
			"userId":    userId,
			"token":     token,
			"platform":  platform,
			"updatedAt": time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (s *NotificationService) DeleteToken(ctx context.Context, userId, token string) error {
	coll := s.DB.Collection(tokenCollectionName)
	_, err := coll.DeleteOne(ctx, bson.M{"userId": userId, "token": token})
	return err
}

func (s *NotificationService) GetTokens(ctx context.Context, userId string) ([]string, error) {
	coll := s.DB.Collection(tokenCollectionName)
	cursor, err := coll.Find(ctx, bson.M{"userId": userId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tokens []string
	for cursor.Next(ctx) {
		var doc struct {
			Token string `bson:"token"`
		}
		if err := cursor.Decode(&doc); err == nil {
			tokens = append(tokens, doc.Token)
		}
	}
	return tokens, nil
}

func (s *NotificationService) FetchTokensByUserIDs(ctx context.Context, userIds []string) (map[string]string, error) {
	coll := s.DB.Collection(tokenCollectionName)
	filter := bson.M{"userId": bson.M{"$in": userIds}}
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result := map[string]string{}
	for cursor.Next(ctx) {
		var doc models.FCMToken
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		result[doc.UserID] = doc.Token
	}
	return result, nil
}

func (s *NotificationService) removeToken(ctx context.Context, userId, token string) {
	coll := s.DB.Collection(tokenCollectionName)
	_, err := coll.DeleteOne(ctx, bson.M{"userId": userId, "token": token})
	if err != nil {
		log.Printf("failed to delete token for user %s: %v", userId, err)
	}
}

func (s *NotificationService) SendNotificationAsync(ctx context.Context, userIds []string, title, body string, data map[string]any) error {
	_, err := s.SendNotificationToUsers(ctx, userIds, title, body, data)
	return err
}

func (s *NotificationService) SendNotificationToUsers(ctx context.Context, userIds []string, title, body string, data map[string]any) (map[string]error, error) {
	userTokenMap, err := s.FetchTokensByUserIDs(ctx, userIds)
	if err != nil {
		return nil, err
	}

	if len(userTokenMap) == 0 {
		return nil, nil
	}

	sem := make(chan struct{}, s.Concurrency)
	var wg sync.WaitGroup
	mu := sync.Mutex{}
	errs := make(map[string]error)

	for userId, token := range userTokenMap {
		wg.Add(1)
		sem <- struct{}{}
		go func(uId, tok string) {
			defer wg.Done()
			defer func() { <-sem }()

			msg := &messaging.Message{
				Token: tok,
				Notification: &messaging.Notification{
					Title: title,
					Body:  body,
				},
				Data: convertData(data),
			}

			_, ferr := s.FCMClient.Send(ctx, msg)
			if ferr != nil {
				job := models.RetryJob{
					UserID:  uId,
					Token:   tok,
					Title:   title,
					Body:    body,
					Data:    data,
					Retries: 1,
				}
				s.enqueueRetry(job)
				mu.Lock()
				errs[uId] = ferr
				mu.Unlock()
				return
			}
		}(userId, token)
	}

	wg.Wait()
	return errs, nil
}

func convertData(data map[string]any) map[string]string {
	result := make(map[string]string)
	for k, v := range data {
		b, _ := json.Marshal(v)
		result[k] = string(b)
	}
	return result
}

func (s *NotificationService) enqueueRetry(job models.RetryJob) {
	if s.Redis == nil {
		return
	}

	job.RetryAt = time.Now().Unix()
	id := fmt.Sprintf("job:%s:%d", job.UserID, job.RetryAt)

	b, _ := json.Marshal(job)

	s.Redis.HSet(context.Background(), id, map[string]string{
		"payload": string(b),
		"retryAt": strconv.FormatInt(job.RetryAt, 10),
	})

	s.Redis.RPush(context.Background(), "notification:retry", id)
}

func (s *NotificationService) StartRetryWorker(ctx context.Context) {
	go func() {
		for {
			result, err := s.Redis.BLPop(ctx, 30*time.Second, "notification:retry").Result()
			if err != nil || len(result) < 2 {
				continue
			}

			jobID := result[1]
			payload, _ := s.Redis.HGet(ctx, jobID, "payload").Result()

			var job models.RetryJob
			if err := json.Unmarshal([]byte(payload), &job); err != nil {
				continue
			}

			msg := &messaging.Message{
				Token: job.Token,
				Notification: &messaging.Notification{
					Title: job.Title,
					Body:  job.Body,
				},
				Data: convertData(job.Data),
			}

			_, ferr := s.FCMClient.Send(ctx, msg)
			if ferr != nil {
				job.Retries++
				if job.Retries < 5 {
					time.Sleep(time.Duration(job.Retries*10) * time.Second)
					s.enqueueRetry(job)
				} else {
					log.Printf("dropping job for user %s after 5 retries", job.UserID)
				}
			} else {
				log.Printf("retried notification delivered to user %s", job.UserID)
			}

			s.Redis.Del(ctx, jobID)
		}
	}()
}

func (s *NotificationService) FetchUserIDsByPincode(ctx context.Context, pincode string) ([]string, error) {
	userColl := s.DB.Collection("farmers")
	cursor, err := userColl.Find(ctx, bson.M{"pincode": pincode})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var userIDs []string
	for cursor.Next(ctx) {
		fmt.Println(cursor.Current)
		var doc struct {
			UserID string `bson:"uid"`
		}
		if err := cursor.Decode(&doc); err == nil {
			userIDs = append(userIDs, doc.UserID)
		}
	}
	return userIDs, nil
}
