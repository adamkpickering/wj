package entry

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const PrettyDateFormat = "January 2, 2006"

var toDoLineRegex *regexp.Regexp
var doneLineRegex *regexp.Regexp
var taskLineRegex *regexp.Regexp

func init() {
	toDoLineRegex = regexp.MustCompile("(?m)^To Do$")
	doneLineRegex = regexp.MustCompile("(?m)^Done$")
	taskLineRegex = regexp.MustCompile("(?m)^([0-9]{2}:[0-9]{2}) ([a-zA-Z0-9,\\_\\-]+) (.*?)$")
}

type Entry struct {
	Date  time.Time
	ToDo  []string
	Done  []string
	Tasks []Task
}

func (entry *Entry) MarshalText() ([]byte, error) {
	buf := &bytes.Buffer{}
	if _, err := fmt.Fprintf(buf, "%s\n\nTo Do\n", entry.Date.Format(PrettyDateFormat)); err != nil {
		return nil, errors.New("failed to write date to buffer")
	}
	for _, toDoItem := range entry.ToDo {
		if _, err := fmt.Fprintf(buf, "- %s\n", toDoItem); err != nil {
			return nil, fmt.Errorf("failed to write to do item %q to buffer", toDoItem)
		}
	}
	if _, err := fmt.Fprintf(buf, "\nDone\n"); err != nil {
		return nil, errors.New("failed to write Done to buffer")
	}
	for _, doneItem := range entry.Done {
		if _, err := fmt.Fprintf(buf, "- %s\n", doneItem); err != nil {
			return nil, fmt.Errorf("failed to write done item %q to buffer", doneItem)
		}
	}
	// Task MarshalText goes here
	return buf.Bytes(), nil
}

func (entry *Entry) UnmarshalText(text []byte) error {
	toDoIndices := toDoLineRegex.FindIndex(text)
	if toDoIndices == nil {
		return errors.New("no match for To Do regex")
	}
	doneIndices := doneLineRegex.FindIndex(text)
	if doneIndices == nil {
		return errors.New("no match for Done regex")
	}
	taskLineIndexPairs := taskLineRegex.FindAllIndex(text, -1)
	if len(taskLineIndexPairs) == 0 {
		return errors.New("no match for task line regex")
	}

	toDoContents := text[toDoIndices[1]:doneIndices[0]]
	entry.ToDo = parseDashList(string(toDoContents))
	doneContents := text[doneIndices[1]:taskLineIndexPairs[0][0]]
	entry.Done = parseDashList(string(doneContents))

	for i := range taskLineIndexPairs {
		var taskContents []byte
		if i+1 == len(taskLineIndexPairs) {
			// the last pair is a special case
			taskContents = text[taskLineIndexPairs[i][0]:]
		} else {
			taskContents = text[taskLineIndexPairs[i][0]:taskLineIndexPairs[i+1][0]]
		}
		task := Task{}
		if err := task.UnmarshalText([]byte(taskContents)); err != nil {
			return fmt.Errorf("failed to parse task: %w", err)
		}
		entry.Tasks = append(entry.Tasks, task)
	}

	// Set the duration of each task. Skip the last task, allowing
	// its duration to remain set to 0.
	for i := 0; i < len(entry.Tasks)-1; i++ {
		duration := entry.Tasks[i+1].StartTime.Sub(entry.Tasks[i].StartTime)
		entry.Tasks[i].Duration = JSONStringDuration(duration)
	}

	return nil
}

func parseDashList(dashListText string) []string {
	lines := strings.Split(strings.TrimSpace(dashListText), "\n")
	dashList := make([]string, 0, len(lines))
	for _, line := range lines {
		dashListEntry := strings.TrimPrefix(line, "- ")
		dashList = append(dashList, dashListEntry)
	}
	return dashList
}
