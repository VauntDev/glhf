package glhf

import (
	"fmt"
	"testing"
)

func TestValidateStatusCode(t *testing.T) {
	testCases := []struct {
		input    int
		expected bool
	}{
		{99, false},
		{100, true},
		{200, true},
		{500, true},
		{999, true},
		{1000, false},
	}

	for v, testCase := range testCases {
		t.Run(fmt.Sprint(v), func(t *testing.T) {

			actualResult := validStatusCode(testCase.input)
			if actualResult != testCase.expected {
				t.Errorf("validateStatusCode(%d) = %t; expected %t", testCase.input, actualResult, testCase.expected)
			}
		})
	}
}
