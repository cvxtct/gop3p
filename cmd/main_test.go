package main

import (
	"testing"
)

func TestParseFiles(t *testing.T) {
	var parseTest = []struct {
		dir      string
		expected []string
	}{
		{dir: "../test_data", expected: []string{"abc.py", "def.py"}},
	}

	for _, tt := range parseTest {
		var p Project
		testname := tt.dir
		t.Run(testname, func(t *testing.T) {
			p.dir = tt.dir
			p.ParseFiles()
			if len(p.files) != len(tt.expected) {
				t.Errorf("got %d, want %d", len(p.files), len(tt.expected))
			}
		})
	}
}
