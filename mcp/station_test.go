package mcp

import "testing"

func TestStation_BuildRRequest(t *testing.T) {
	station := NewLocalStation()
	request := station.BuildReadRequest("D", 300, 3)

	if request != "500000FFFF03000C001000010400002C0100A80300" {
		t.Fatalf("expected %v but actual is %v", "500000FFFF03000C001000010400002C0100A80300", request)
	}

	request2 := station.BuildReadRequest("D", 500, 50)
	if request2 != "500000FFFF03000C00100001040000F40100A83200" {
		t.Fatalf("expected %v but actual is %v", "500000FFFF03000C00100001040000F40100A83200", request2)
	}
}
