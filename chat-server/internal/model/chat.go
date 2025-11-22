package model

import "time"

type Chat struct {
	ID        int64
	Usernames []string
}

type Message struct {
	From      string
	Text      string
	Timestamp time.Time
}
