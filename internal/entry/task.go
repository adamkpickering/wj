package entry

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var taskRegex *regexp.Regexp

func init() {
	taskRegex = regexp.MustCompile("(?s)([0-9]{2}:[0-9]{2}) ([a-zA-Z0-9,\\_\\-]+) (.*?)\n(.*)$")
}

type JSONStringDuration time.Duration

func (d JSONStringDuration) MarshalJSON() ([]byte, error) {
	asDuration := time.Duration(d)
	output := fmt.Sprintf("%q", asDuration)
	return []byte(output), nil
}

type Task struct {
	Title     string
	Duration  JSONStringDuration
	StartTime time.Time
	Tags      []string
	Body      string
}

func (task *Task) UnmarshalText(text []byte) error {
	submatches := taskRegex.FindSubmatch(text)
	if submatches == nil {
		return fmt.Errorf("no match for task regex for text %q", bytes.TrimSpace(text)[:40])
	}

	rawStartTime := string(submatches[1])
	parsedTime, err := time.Parse("15:04", rawStartTime)
	if err != nil {
		return fmt.Errorf("failed to parse time %q: %w", rawStartTime, err)
	}
	task.StartTime = parsedTime

	task.Tags = strings.Split(string(submatches[2]), ",")
	task.Title = string(submatches[3])
	task.Body = strings.TrimSpace(string(submatches[4]))

	return nil
}
