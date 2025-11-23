package model

import "time"

type UserLog struct {
	ID        int64
	Action    string
	EntityID  int64
	CreatedAt time.Time
}
