/*
Copyright © 2024 Chris Greaves cjgreaves97@hotmail.co.uk

See the file COPYING in the root of this repository for details.
*/

package processor

import (
	"testing"
)

func TestCreateFuncsFromModel(t *testing.T) {
	model := Model{
		Var: map[string]string{"key1": "value1"},
		Env: map[string]string{"key2": "value2"},
	}

	funcs := CreateFuncsFromModel(model)

	tests := []struct {
		name     string
		funcName string
		key      string
		expected string
	}{
		{
			name:     "var function",
			funcName: "var",
			key:      "key1",
			expected: "value1",
		},
		{
			name:     "env function",
			funcName: "env",
			key:      "key2",
			expected: "value2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn, ok := funcs[tt.funcName].(func(string) string)
			if !ok {
				t.Fatalf("unexpected type assertion error")
			}

			if got := fn(tt.key); got != tt.expected {
				t.Errorf("expected '%s', but got '%s'", tt.expected, got)
			}
		})
	}
}
