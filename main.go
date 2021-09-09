package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
	"virtualhostingserver/common"
)

const (
	HTTP_HEADER_STATUS_OK        = "HTTP/1.1 200 OK\n\r\n\r"
	HTTP_HEADER_STATUS_NOT_FOUND = "HTTP/1.1 404 Not Found\n\r\n\r"
)

var GlobalConfig common.Config

func init() {
	var err error
	GlobalConfig, err = common.ReadConfig("config.yml")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	listenAddr := fmt.Sprintf("%s:%d", GlobalConfig.IP, GlobalConfig.Port)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		conn.SetDeadline(time.Now().Add(time.Second * GlobalConfig.Timeout))

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	httpHeaders := getHttpHeaders(conn)

	template, headerStatus := serveDirectory(httpHeaders)
	conn.Write([]byte(headerStatus + template + "\n\r"))
}

// TODO: Move all below this to a 'package parsers'
func getHttpHeaders(conn net.Conn) HttpHeaders {
	reader := bufio.NewReader(conn)
	var results []string
	areMoreLinesToRead := true
	for areMoreLinesToRead {
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Fatal("Couldn't read from connection")
		}

		areMoreLinesToRead = strings.IndexByte(string(line), ' ') != -1
		results = append(results, string(line))
	}

	toReturn := HttpHeaders{}
	err := toReturn.fromStringArray(results)
	if err != nil {
		log.Fatal(err)
	}

	return toReturn
}
