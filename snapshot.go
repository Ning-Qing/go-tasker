package gotasker

import "time"

type Snapshot struct {
	ID        string
	TaskID    string
	Step      uint8
	Input     []byte
	Data      []byte
	CreatedAt time.Time
}
