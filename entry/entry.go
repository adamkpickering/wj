package entry

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type ParseState string

var taskLineRegex regexp.Regexp

func init() {
	taskLineRegex = *regexp.MustCompile(`^([0-9]{2}:[0-9]{2}) ([a-z0-9,]+) (.*)$`)
}

const (
	dateParseState  ParseState = "date"
	doneParseState  ParseState = "done"
	toDoParseState  ParseState = "toDo"
	tasksParseState ParseState = "tasks"
)

const PrettyDateFormat = "January 2, 2006"

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
	var (
		task Task
		err  error
	)
	contents := string(text)
	taskContentLines := make([]string, 0, 100)
	lines := strings.Split(contents, "\n")
	var parseState ParseState = dateParseState
	for _, line := range lines {
		if parseState == dateParseState {
			// Looking for entry date and To Do line
			if line == "To Do" {
				if entry.Date.IsZero() {
					return errors.New("no date in entry")
				}
				parseState = toDoParseState
			}
			date, err := time.Parse(PrettyDateFormat, line)
			if err == nil {
				entry.Date = date
			}
		} else if parseState == toDoParseState {
			// Parsing to do statements and skipping empty lines
			if line == "Done" {
				parseState = doneParseState
				continue
			} else if strings.HasPrefix(line, "- ") {
				toDoText := strings.TrimPrefix(line, "- ")
				entry.ToDo = append(entry.ToDo, toDoText)
			}
		} else if parseState == doneParseState {
			// Parsing done statements and skipping empty lines
			if taskLineRegex.MatchString(line) {
				parseState = tasksParseState

				task, err = partialTaskFromTitleLine(line)
				if err != nil {
					return fmt.Errorf("failed to parse first title line %q: %w", line, err)
				}
			} else if strings.HasPrefix(line, "- ") {
				doneText := strings.TrimPrefix(line, "- ")
				entry.Done = append(entry.Done, doneText)
			}
		} else if parseState == tasksParseState {
			if taskLineRegex.MatchString(line) {
				newTask, err := partialTaskFromTitleLine(line)
				if err != nil {
					return fmt.Errorf("failed to parse title line %q: %w", line, err)
				}

				// set .Content of previous task
				task.Content = strings.Join(taskContentLines, "\n")
				taskContentLines = make([]string, 0, 100)

				// calculate duration of previous task
				task.Duration = newTask.StartTime.Sub(task.StartTime)

				entry.Tasks = append(entry.Tasks, task)

				task = newTask
			} else {
				taskContentLines = append(taskContentLines, line)
			}
		}
	}

	// set .Content of last task and add it to list
	task.Content = strings.Join(taskContentLines, "\n")
	entry.Tasks = append(entry.Tasks, task)

	return nil
}

// Constructs a Task and reads relevant fields from a Task title line
// into that Task. If the title line does not contain info on
// a given field of the Task, that field is left as its zero value.
func partialTaskFromTitleLine(line string) (Task, error) {
	task := Task{}

	result := taskLineRegex.FindStringSubmatch(line)
	if len(result) != 4 {
		return Task{}, fmt.Errorf("failed to parse line %q", line)
	}

	// parse start time
	rawStartTime := result[1]
	parsedTime, err := time.Parse("15:04", rawStartTime)
	if err != nil {
		return Task{}, fmt.Errorf("failed to parse time %q: %w", rawStartTime, err)
	}
	task.StartTime = parsedTime

	// parse tags
	rawTags := result[2]
	parts := strings.Split(rawTags, ",")
	task.Tags = append(task.Tags, parts...)

	task.Title = result[3]

	return task, nil
}
