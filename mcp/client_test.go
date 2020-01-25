package mcp

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestClient3E_Read(t *testing.T) {
	// running only when there is and plc that can be accepted mc protocol

	client, err := New3EClient("10.23.3.117", 1280, NewLocalStation())
	if err != nil {
		t.Fatalf("PLC does not exists? %v", err)
	}

	// 1 device
	resp1, err := client.Read("D", 100, 1)
	if err != nil {
		t.Fatalf("unexpected mcp read err: %v", err)
	}

	if len(resp1) != 13 {
		t.Fatalf("expected %v but actual is %v", 13, len(resp1))
	}
	if hex.EncodeToString(resp1) != strings.ReplaceAll("d000 00 ff ff03 0004 0000 0000 00", " ", "") {
		t.Fatalf("expected %v but actual is %v", "d00000ffff0300040000000000", hex.EncodeToString(resp1))
	}

	// 3 device
	resp2, err := client.Read("D", 100, 5)
	if err != nil {
		t.Fatalf("unexpected mcp read err: %v", err)
	}

	if len(resp2) != 21 {
		t.Fatalf("expected %v but actual is %v", 21, len(resp2))
	}

	if hex.EncodeToString(resp2) != strings.ReplaceAll("d000 00 ff ff03 000c 0000 0000 000000000000000000", " ", "") {
		t.Fatalf("expected %v but actual is %v", "d00000ffff03000c00000000000000000000000000", hex.EncodeToString(resp2))
	}

}

func TestClient3E_Ping(t *testing.T) {
	// running only when there is and plc that can be accepted mc protocol

	client, err := New3EClient("10.23.3.117", 1280, NewLocalStation())
	if err != nil {
		t.Fatalf("PLC does not exists? %v", err)
	}

	if err := client.HealthCheck(); err != nil {
		t.Fatalf("unexpected error occured %v", err)
	}
}
