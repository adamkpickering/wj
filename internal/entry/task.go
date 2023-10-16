package entry

import (
	"time"
)

type Task struct {
	StartTime time.Time
	Tags      []string
	Title     string
	Content   string
	Duration  time.Duration
}
