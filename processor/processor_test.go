package processor

import (
	"reflect"
	"testing"
)

func TestConvertVarsArrayToMap(t *testing.T) {
	testCases := []struct {
		desc   string
		input  []string
		output map[string]string
	}{
		{
			desc:   "will convert valid vars",
			input:  []string{"foo=bar"},
			output: map[string]string{"foo": "bar"},
		},
		{
			desc:   "will convert with more than one =",
			input:  []string{"foo=bar=="},
			output: map[string]string{"foo": "bar=="},
		},
		{
			desc:   "will ignore vars without an =",
			input:  []string{"foo=bar", "boo"},
			output: map[string]string{"foo": "bar"},
		},
	}
	for _, c := range testCases {
		t.Run(c.desc, func(t *testing.T) {
			result := convertVarsArrayToMap(c.input)

			if eq := reflect.DeepEqual(result, c.output); !eq {
				t.Errorf("got '%v', expected '%v'", result, c.output)
			}
		})
	}
}
