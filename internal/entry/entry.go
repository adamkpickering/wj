package entry

import (
	"time"
)

type Entry struct {
	ToDo  []string
	Done  []string
	Tasks []Task
}

type Task struct {
	StartTime time.Time
	Tags      []string
	Title     string
	Content   string
	Duration  time.Duration
}

// type JSONStringDuration time.Duration

// func (d JSONStringDuration) MarshalJSON() ([]byte, error) {
// 	asDuration := time.Duration(d)
// 	output := fmt.Sprintf("%q", asDuration)
// 	return []byte(output), nil
// }
