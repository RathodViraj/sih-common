package models

import "time"

type Error struct {
	Data    any       `json:"data"`
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}
