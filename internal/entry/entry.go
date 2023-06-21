package entry

import (
	"fmt"
	"time"
)

type Entry struct {
	ToDo  []string
	Done  []string
	Tasks []Task
}

type Task struct {
	Description string
	Duration    JSONStringDuration
	StartTime   time.Time
	Content     string
}

type JSONStringDuration time.Duration

func (rawDuration JSONStringDuration) Pretty() (prettyDuration string) {
	duration := time.Duration(rawDuration)
	hours := duration / time.Hour
	minutes := (duration - hours*time.Hour) / time.Minute
	if hours > 0 {
		prettyDuration = fmt.Sprintf("%dh%dm", hours, minutes)
	} else {
		prettyDuration = fmt.Sprintf("%dm", minutes)
	}
	return
}

func (d JSONStringDuration) MarshalJSON() ([]byte, error) {
	asDuration := time.Duration(d)
	output := fmt.Sprintf("%q", asDuration)
	return []byte(output), nil
}
