package models

import "time"

type Success struct {
	Response any       `json:"response"`
	Code     int       `json:"code"`
	Message  string    `json:"message"`
	Time     time.Time `json:"time"`
}
