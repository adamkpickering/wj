package cmd

import (
	"testing"
	"time"
)

func TestPretty(t *testing.T) {
	for _, testCase := range []struct {
		Input  time.Duration
		Output string
	}{
		{Input: 9 * time.Minute, Output: "    9m"},
		{Input: 52 * time.Minute, Output: "   52m"},
		{Input: 4*time.Hour + 9*time.Minute, Output: " 4h 9m"},
		{Input: 13*time.Hour + 9*time.Minute, Output: "13h 9m"},
		{Input: 4*time.Hour + 29*time.Minute, Output: " 4h29m"},
		{Input: 13*time.Hour + 29*time.Minute, Output: "13h29m"},
	} {
		t.Run("should produce the correct output", func(t *testing.T) {
			if output := pretty(testCase.Input); output != testCase.Output {
				t.Errorf("input %v output %q (expected %q)", testCase.Input, output, testCase.Output)
			}
		})
	}
}
