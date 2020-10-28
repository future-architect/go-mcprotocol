package mcp

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

type Client interface {
	Read(deviceName string, offset, numPoints int64) ([]byte, error)
	Write(deviceName string, offset, numPoints int64, writeData []byte) ([]byte, error)
	HealthCheck() error
	Connect() error
	Close() error
}

// client3E is 3E frame mcp client
type client3E struct {
	// PLC address
	tcpAddr *net.TCPAddr
	// PLC station
	stn *station

	// Connect & Read Write timeout
	Timeout time.Duration

	// TCP connection
	mu   sync.Mutex
	conn net.Conn
}

func New3EClient(host string, port int, stn *station) (Client, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		return nil, err
	}
	return &client3E{tcpAddr: tcpAddr, stn: stn}, nil
}

// MELSECコミュニケーションプロトコル p180
// 11.4折返しテスト
func (c *client3E) HealthCheck() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	requestStr := c.stn.BuildHealthCheckRequest()

	// binary protocol
	payload, err := hex.DecodeString(requestStr)
	if err != nil {
		return err
	}

	// Connection established if not connect
	if err = c.connect(); err != nil {
		return err
	}

	// Set write and read timeout if set timeout
	if c.Timeout > 0 {
		deadLine := time.Now().Add(c.Timeout)
		if err = c.conn.SetDeadline(deadLine); err != nil {
			return err
		}
	}

	// Send message
	if _, err = c.conn.Write(payload); err != nil {
		return err
	}

	// Receive message
	readBuff := make([]byte, 30)
	readLen, err := c.conn.Read(readBuff)
	if err != nil {
		return err
	}

	resp := readBuff[:readLen]

	if readLen != 18 {
		return errors.New("plc connect test is fail: return length is [" + fmt.Sprintf("%X", resp) + "]")
	}

	// decodeString is 折返しデータ数ヘッダ[1byte]
	if "0500" != fmt.Sprintf("%X", resp[11:13]) {
		return errors.New("plc connect test is fail: return header is [" + fmt.Sprintf("%X", resp[11:13]) + "]")
	}

	//  折返しデータ[5byte]=ABCDE
	if "4142434445" != fmt.Sprintf("%X", resp[13:18]) {
		return errors.New("plc connect test is fail: return body is [" + fmt.Sprintf("%X", resp[13:18]) + "]")
	}

	return nil
}

// Read is send read command to remote plc by mc protocol
// deviceName is device code name like 'D' register.
// offset is device offset addr.
// numPoints is number of read device points.
func (c *client3E) Read(deviceName string, offset, numPoints int64) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	requestStr := c.stn.BuildReadRequest(deviceName, offset, numPoints)

	// TODO binary protocol
	payload, err := hex.DecodeString(requestStr)
	if err != nil {
		return nil, err
	}

	// Connection established if not connect
	if err = c.connect(); err != nil {
		return nil, err
	}

	// Set write and read timeout if set timeout
	if c.Timeout > 0 {
		deadLine := time.Now().Add(c.Timeout)
		if err = c.conn.SetDeadline(deadLine); err != nil {
			return nil, err
		}
	}

	// Send message
	if _, err = c.conn.Write(payload); err != nil {
		return nil, err
	}

	// Receive message
	readBuff := make([]byte, 22+2*numPoints) // 22 is response header size. [sub header + network num + unit i/o num + unit station num + response length + response code]
	readLen, err := c.conn.Read(readBuff)
	if err != nil {
		return nil, err
	}

	return readBuff[:readLen], nil
}

// Write is send write command to remote plc by mc protocol
// deviceName is device code name like 'D' register.
// offset is device offset addr.
// writeData is data to write.
// numPoints is number of write device points.
// writeData is the data to be written. If writeData is larger than 2*numPoints bytes,
// data larger than 2*numPoints bytes is ignored.
func (c *client3E) Write(deviceName string, offset, numPoints int64, writeData []byte) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	requestStr := c.stn.BuildWriteRequest(deviceName, offset, numPoints, writeData)
	payload, err := hex.DecodeString(requestStr)
	if err != nil {
		return nil, err
	}

	// Connection established if not connect
	if err = c.connect(); err != nil {
		return nil, err
	}

	// Set write and read timeout if set timeout
	if c.Timeout > 0 {
		deadLine := time.Now().Add(c.Timeout)
		if err = c.conn.SetDeadline(deadLine); err != nil {
			return nil, err
		}
	}

	// Send message
	if _, err = c.conn.Write(payload); err != nil {
		return nil, err
	}

	// Receive message
	readBuff := make([]byte, 22) // 22 is response header size. [sub header + network num + unit i/o num + unit station num + response length + response code]

	readLen, err := c.conn.Read(readBuff)
	if err != nil {
		return nil, err
	}

	return readBuff[:readLen], nil
}

// Close closes the connection. Close only if the connection exists.
func (c *client3E) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.close()
}

// Connect establishes a new connection to the tcp Address.
// Establishes a connection only if a connection has not been established.
// If a connection already exists, no action is taken.
func (c *client3E) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.connect()
}

func (c *client3E) connect() error {
	if c.conn == nil {
		dialer := net.Dialer{Timeout: c.Timeout}
		conn, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", c.tcpAddr.IP.String(), c.tcpAddr.Port))
		if err != nil {
			return err
		}
		c.conn = conn
	}
	return nil
}

func (c *client3E) close() error {
	var err error
	if c.conn != nil {
		err = c.conn.Close()
		c.conn = nil
	}
	return err
}
