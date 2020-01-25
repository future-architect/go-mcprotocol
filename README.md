# go-mcprotocol

go-mcprotocol is a library for PLC (Programmable Logic Controller) access

## Project Status

**Work In Progress**

## Usage for Library

You can read plc register bellow codes.

```go
	client, _ := mcp.New3EClient(opts.Host, opts.Port, mcp.NewLocalStation())
	read, _ := client.Read("D", 100, 3)
	registerBinary, _ := mcp.NewParser().Do(read)

	fmt.Println(string(registerBinary.Payload))
```

#### Health Check

```go
	if err := client.HealthCheck(); err != nil {
		log.Fatalf("failed health check for plc: %v", err)
	}
```

## Usage Tool

## Output file format

Format is CSV. Items are timestamp and Base64 encoded MC Protocol response.

```csv
2019-10-07T07:08:00.3623052Z,0AAA//8DAAwAAAAAAAAAAAAAAAAA
2019-10-07T07:08:00.8622182Z,0AAA//8DAAwAAAAAAAAAAAAAAAAA
2019-10-07T07:08:01.3616205Z,0AAA//8DAAwAAAAAAAAAAAAAAAAA
...
```

## Usage for tool

Collect the register values​of PLC by MC Protocol(MELSEC Communication Protocol).
This tools is gather plc register data and dump local files.

### Examples

```bash
$ plcmirror -H <Your PLC Host> -P <Your PLC Port> --device D --offset 100 --num 10 --dir /var/log/plcmirror
```

### Options

```bash
> plcmirror -help
  Usage:
    plcmirror [OPTIONS]
  
  Application Options:
    /H, /host:      PLC hostname
    /P, /port:      Melsec communication protocol port number
    /D, /device:    Register name like D that is mirror target
    /O, /offset:    PLC register offset addr that is mirror target
    /N, /num:       number of device points
        /dir:       file output path (default: .)
    /I, /interval:  mirroring interval [milli sec] (default: 500)
  
  Help Options:
    /?              Show this help message
    /h, /help       Show this help message
```



# License
Apache 2
