package models

type RetryJob struct {
	UserID  string         `json:"userId"`
	Token   string         `json:"token"`
	Title   string         `json:"title"`
	Body    string         `json:"body"`
	Data    map[string]any `json:"data"`
	Retries int            `json:"retries"`
	RetryAt int64          `json:"retryAt"`
}
