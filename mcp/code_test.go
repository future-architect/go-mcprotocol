package mcp

import (
	"encoding/hex"
	"testing"
)

func TestCode_EncodeHex(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{input: "0401", expected: "0104"},
	}

	for _, v := range cases {
		actual, err := Binary.EncodeHex(v.input)
		if err != nil {
			t.Errorf("something wrong: input is %v", v.input)
			continue
		}

		if hex.EncodeToString(actual) != v.expected {
			t.Errorf("wrong result: expected is %v but actual is %v", "0104", hex.EncodeToString(actual))
		}
	}
}
