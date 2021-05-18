package mcp

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	SUB_HEADER = "5000" // 3Eフレームでは固定

	HEALTH_CHECK_COMMAND    = "1906" // binary mode expression. if ascii mode then 0619
	HEALTH_CHECK_SUBCOMMAND = "0000"

	READ_COMMAND         = "0104" // binary mode expression. if ascii mode then 0401
	READ_SUB_COMMAND     = "0000"
	BIT_READ_SUB_COMMAND = "0100"

	WRITE_COMMAND     = "0114" // binary mode expression. if ascii mode then 1401
	WRITE_SUB_COMMAND = "0000"

	MONITORING_TIMER = "1000" // 3[sec]
)

// deviceCodes is device name and hex value map
var deviceCodes = map[string]string{
	"X": "9C",
	"Y": "9D",
	"M": "90",
	"L": "92",
	"F": "93",
	"V": "94",
	"B": "A0",
	"W": "B4",
	"D": "A8",
}

// Each single PLC that is connected on MELSECNET and CC-Link IE is called a station.
type station struct {
	// PLC Network number
	networkNum string
	// PC Number
	pcNum string
	// PLC stn Unit I/O Number
	unitIONum string
	// PLC stn Unit Station Number
	unitStationNum string
}

func NewStation(networkNum, pcNum, unitIONum, unitStationNum string) *station {
	return &station{
		networkNum:     networkNum,
		pcNum:          pcNum,
		unitIONum:      unitIONum,
		unitStationNum: unitStationNum,
	}
}

// local stn stn. local stn is 自局.
func NewLocalStation() *station {
	return &station{
		networkNum:     "00",   // 自局の場合は00固定
		pcNum:          "FF",   // 自局の場合はFF固定
		unitIONum:      "FF03", // マルチドロップ接続などでない場合はFF03固定値
		unitStationNum: "00",   // マルチドロップ接続などでない場合は00固定値
	}
}

func (h *station) BuildHealthCheckRequest() string {

	returnDataNum := "0500"    // 5 device. if ascii mode then 0005
	returnData := "4142434445" // value is "ABCDE".

	requestStr := HEALTH_CHECK_COMMAND + HEALTH_CHECK_SUBCOMMAND + returnDataNum + returnData

	// data length
	requestCharLen := len(MONITORING_TIMER+requestStr) / 2 // 1byte=2char
	dataLenBuff := new(bytes.Buffer)
	_ = binary.Write(dataLenBuff, binary.LittleEndian, int64(requestCharLen))
	dataLen := fmt.Sprintf("%X", dataLenBuff.Bytes()[0:2]) // 2byte固定

	return SUB_HEADER +
		h.networkNum +
		h.pcNum +
		h.unitIONum +
		h.unitStationNum +
		dataLen +
		MONITORING_TIMER +
		requestStr
}

// BuildReadRequest represents MCP read as word command.
// deviceName is device code name like 'D' register.
// offset is device offset addr.
// numPoints is number of read device points.
func (h *station) BuildReadRequest(deviceName string, offset, numPoints int64) string {

	// get device symbol hex layout
	deviceCode := deviceCodes[deviceName]

	// offset convert to little endian layout
	// MELSECコミュニケーションプロトコル リファレンス(p67) MELSEC-Q/L: 3[byte], MELSEC iQ-R: 4[byte]
	offsetBuff := new(bytes.Buffer)
	_ = binary.Write(offsetBuff, binary.LittleEndian, offset)
	offsetHex := fmt.Sprintf("%X", offsetBuff.Bytes()[0:3]) // 仮にQシリーズとするので3byte trim

	// read points
	pointsBuff := new(bytes.Buffer)
	_ = binary.Write(pointsBuff, binary.LittleEndian, numPoints)
	points := fmt.Sprintf("%X", pointsBuff.Bytes()[0:2]) // 2byte固定

	// data length
	requestCharLen := len(MONITORING_TIMER+READ_COMMAND+READ_SUB_COMMAND+deviceCode+offsetHex+points) / 2 // 1byte=2char
	dataLenBuff := new(bytes.Buffer)
	_ = binary.Write(dataLenBuff, binary.LittleEndian, int64(requestCharLen))
	dataLen := fmt.Sprintf("%X", dataLenBuff.Bytes()[0:2]) // 2byte固定

	return SUB_HEADER +
		h.networkNum +
		h.pcNum +
		h.unitIONum +
		h.unitStationNum +
		dataLen +
		MONITORING_TIMER +
		READ_COMMAND +
		READ_SUB_COMMAND +
		offsetHex +
		deviceCode +
		points
}

// BuildReadRequest represents MCP read as bit command.
// deviceName is device code name like 'D' register.
// offset is device offset addr.
// numPoints is number of read device points.
func (h *station) BuildBitReadRequest(deviceName string, offset, numPoints int64) string {

	// get device symbol hex layout
	deviceCode := deviceCodes[deviceName]

	// offset convert to little endian layout
	// MELSECコミュニケーションプロトコル リファレンス(p67) MELSEC-Q/L: 3[byte], MELSEC iQ-R: 4[byte]
	offsetBuff := new(bytes.Buffer)
	_ = binary.Write(offsetBuff, binary.LittleEndian, offset)
	offsetHex := fmt.Sprintf("%X", offsetBuff.Bytes()[0:3]) // 仮にQシリーズとするので3byte trim

	// read points
	pointsBuff := new(bytes.Buffer)
	_ = binary.Write(pointsBuff, binary.LittleEndian, numPoints)
	points := fmt.Sprintf("%X", pointsBuff.Bytes()[0:2]) // 2byte固定

	// data length
	requestCharLen := len(MONITORING_TIMER+READ_COMMAND+BIT_READ_SUB_COMMAND+deviceCode+offsetHex+points) / 2 // 1byte=2char
	dataLenBuff := new(bytes.Buffer)
	_ = binary.Write(dataLenBuff, binary.LittleEndian, int64(requestCharLen))
	dataLen := fmt.Sprintf("%X", dataLenBuff.Bytes()[0:2]) // 2byte固定

	return SUB_HEADER +
		h.networkNum +
		h.pcNum +
		h.unitIONum +
		h.unitStationNum +
		dataLen +
		MONITORING_TIMER +
		READ_COMMAND +
		BIT_READ_SUB_COMMAND +
		offsetHex +
		deviceCode +
		points
}

// BuildWriteRequest represents MCP write command.
// deviceName is device code name like 'D' register.
// offset is device offset addr.
// writeData is data to write.
// numPoints is number of write device points.
// writeData is the data to be written. If writeData is larger than 2*numPoints bytes,
// data larger than 2*numPoints bytes is ignored.
func (h *station) BuildWriteRequest(deviceName string, offset, numPoints int64, writeData []byte) string {

	// get device symbol hex layout
	deviceCode := deviceCodes[deviceName]

	// offset convert to little endian layout
	// MELSECコミュニケーションプロトコル リファレンス(p67) MELSEC-Q/L: 3[byte], MELSEC iQ-R: 4[byte]
	offsetBuff := new(bytes.Buffer)
	_ = binary.Write(offsetBuff, binary.LittleEndian, offset)
	offsetHex := fmt.Sprintf("%X", offsetBuff.Bytes()[0:3]) // 仮にQシリーズとするので3byte trim

	// convert write data to little endian word
	writeBuff := new(bytes.Buffer)
	_ = binary.Write(writeBuff, binary.LittleEndian, writeData)
	writeHex := fmt.Sprintf("%X", writeBuff.Bytes()[0:2*numPoints]) // 2 byte per 1 device point

	// write points
	pointsBuff := new(bytes.Buffer)
	_ = binary.Write(pointsBuff, binary.LittleEndian, numPoints)
	points := fmt.Sprintf("%X", pointsBuff.Bytes()[0:2]) // 2byte固定

	// data length
	requestCharLen := len(MONITORING_TIMER+WRITE_COMMAND+WRITE_SUB_COMMAND+deviceCode+offsetHex+points+writeHex) / 2 // 1byte=2char
	dataLenBuff := new(bytes.Buffer)
	_ = binary.Write(dataLenBuff, binary.LittleEndian, int64(requestCharLen))
	dataLen := fmt.Sprintf("%X", dataLenBuff.Bytes()[0:2]) // 2byte固定
	return SUB_HEADER +
		h.networkNum +
		h.pcNum +
		h.unitIONum +
		h.unitStationNum +
		dataLen +
		MONITORING_TIMER +
		WRITE_COMMAND +
		WRITE_SUB_COMMAND +
		offsetHex +
		deviceCode +
		points +
		writeHex
}

func (h *station) BuildAccessPath() {

}
