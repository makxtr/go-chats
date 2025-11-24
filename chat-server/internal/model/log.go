package model

import "time"

type ChatLog struct {
	ID        int64
	Action    string
	EntityID  int64
	CreatedAt time.Time
}
