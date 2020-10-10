package mcp

import (
	"encoding/hex"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestParser_Do(t *testing.T) {
	mcResp, _ := hex.DecodeString("d00000ffff0300040000000000")

	p := NewParser()
	response, err := p.Do(mcResp)
	if err != nil {
		t.Fatalf("unexpected parser err: %v", err)
	}

	expected := &Response{
		SubHeader:      "D000",
		NetworkNum:     "00",
		PCNum:          "FF",
		UnitIONum:      "FF03",
		UnitStationNum: "00",
		DataLen:        "0400",
		EndCode:        "0000",
		Payload:        []uint8{0x00, 0x00},
		ErrInfo:        nil,
	}

	if diff := cmp.Diff(response, expected); diff != "" {
		t.Errorf("parse Resp differs: (-got +want)\n%s", diff)
	}
}
