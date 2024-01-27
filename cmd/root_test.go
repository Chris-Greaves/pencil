package cmd

import "testing"

func TestValidateVars(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		want  error
	}{
		{"no errors when valid var", "foo=bar", nil},
		{"no error when value contains =", "foo=bar==", nil},
		{"error when input doesn't contain =", "foo", ErrInvalidVariable},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vars := []string{
				tt.input,
			}

			output, err := validateDirectVars(vars)
			t.Logf("The output was %v", output)

			if err != tt.want {
				t.Errorf("got '%s', wanted '%s'", err, tt.want)
			}
		})
	}
}
