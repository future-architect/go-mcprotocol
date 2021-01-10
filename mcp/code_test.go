package mcp

import (
	"encoding/hex"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCode_EncodeHex(t *testing.T) {
	tests := []struct {
		name      string
		c         Code
		inputCode string
		want      string
		wantErr   bool
	}{
		{name: "Binary mode", c: Binary, inputCode: "0401", want: "0104", wantErr: false},
		{name: "Binary mode max value", c: Binary, inputCode: "FFFF", want: "ffff", wantErr: false},
		{name: "Ascii mode", c: Ascii, inputCode: "0401", want: "0401", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.EncodeHex(tt.inputCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			var gotStr string
			switch tt.c {
			case Ascii:
				gotStr = string(got)
			case Binary:
				gotStr = hex.EncodeToString(got)
			default:
				t.Errorf("not supported Code, Code is %v", tt.c)
				return
			}

			if diff := cmp.Diff(tt.want, gotStr); diff != "" {
				t.Errorf("EncodeHex() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
