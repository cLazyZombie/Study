package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Go simple http server started")
	fmt.Printf("0x0d0a = %s", string([]byte{0x0d, 0x0a}))
	fmt.Printf("\\n = %s", hex.EncodeToString([]byte("\n")))

	ln, err := net.Listen("tcp", ":8787")
	if err != nil {
		fmt.Printf("error on Listen. %s\n", err.Error())
		os.Exit(1)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("error on Accept(). %s\n", err.Error())
		} else {
			go handleConnection(conn)
		}
	}
}

func handleConnection(conn net.Conn) {
	fmt.Printf("handleConnection started\n")

	remoteAddr := conn.RemoteAddr()
	fmt.Printf("RemoteAddr: %s\n", remoteAddr.String())

	reader := bufio.NewReader(conn)

	// read request

	command := ""
	commandArg := ""
	for {
		readStr, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Read Error. %s", err.Error())
				conn.Close()
				return
			} else {
				fmt.Printf("Read Finished")
				break
			}
		} else {
			fmt.Printf("Read Content: %s", readStr)

			fmt.Printf("Hex: %s\n", hex.EncodeToString([]byte(readStr)))

			// read finished
			if len(readStr) == 0 || readStr == string([]byte{0x0d, 0x0a}) {
				break
			}

			tokens := strings.Split(readStr, " ")
			if len(tokens) > 2 && tokens[0] == "GET" {
				command = "GET"
				commandArg = tokens[1]
			}
		}
	}

	switch command {
	case "GET":
		if commandArg == "/" {

			s := `HTTP/1.1 200 OK
Date: Mon, 23 May 2005 22:38:34 GMT
Content-Type: text/html; charset=UTF-8
Content-Encoding: UTF-8
Content-Length: 138
Last-Modified: Wed, 08 Jan 2003 23:11:55 GMT
Server: Apache/1.3.3.7 (Unix) (Red-Hat/Linux)
ETag: "3f80f-1b6-3e1cb03b"
Accept-Ranges: bytes
Connection: close

<html>
<head>
  <title>An Example Page</title>
</head>
<body>
  Hello World, this is a very simple HTML document.
</body>
</html>
`
			conn.Write([]byte(s))
			conn.Write([]byte{0x0d, 0x0a})
			conn.Close()
		}
	}
}
