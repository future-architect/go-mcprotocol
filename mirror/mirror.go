package mirror

import (
	"encoding/base64"
	"fmt"
	"github.com/pj-cancan/plcmirror/mcp"
	"io"
	"log"
	"sync"
	"time"
)

var (
	mu       sync.Mutex
	permitRW bool
)

type fileMirror struct {
	c          mcp.Client
	w          io.Writer
	deviceName string
	offset     int64
	numPoints  int64
	interval   time.Duration
}

func NewFileMirror(c mcp.Client, w io.Writer, deviceName string, offset, numPoints int64, interval time.Duration) *fileMirror {
	return &fileMirror{
		c:          c,
		w:          w,
		deviceName: deviceName,
		offset:     offset,
		numPoints:  numPoints,
		interval:   interval,
	}
}

func (m *fileMirror) RunAndServe() error {
	c := time.Tick(m.interval)
	for {
		select {
		case <-c:
			go m.readAndWrite()
		}
	}
}

func (m *fileMirror) readAndWrite() {

	if !m.Lock() {
		log.Printf("[INFO] skip readAndWrite because goroutine cannnot get lock")
		// skip because cannot get lock
		return
	}
	defer m.Unlock()

	bytes, err := m.c.Read(m.deviceName, m.offset, m.numPoints)
	if err != nil {
		log.Printf("[ERROR] plc read error: %v\n", err)
		return
	}
	payload := base64.StdEncoding.EncodeToString(bytes)
	now := time.Now().UTC().Format(time.RFC3339Nano)

	_, err = m.w.Write([]byte(fmt.Sprintf("%v,%v\n", now, payload)))
	if err != nil {
		log.Printf("[ERROR] plc data file write error: %v\n", err)
	}

}

// Guards duplicate plc access and skip when delay read operation for reducing plc workload
func (m *fileMirror) Lock() bool {
	mu.Lock()
	defer mu.Unlock()
	if permitRW {
		return false
	}
	permitRW = true
	return true
}

func (m *fileMirror) Unlock() {
	mu.Lock()
	permitRW = false
	mu.Unlock()
}
