package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"time"
)

const (
	HTTP_HEADER_STATUS_OK        = "HTTP/1.1 200 OK\n\r\n\r"
	HTTP_HEADER_STATUS_NOT_FOUND = "HTTP/1.1 404 Not Found\n\r\n\r"
	CONFIG_PATH                  = "config.json"
)

var GlobalConfig Config

func main() {
	GlobalConfig = loadConfig()
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		conn.SetDeadline(time.Now().Add(time.Second * 60)) // Set timeout to 1 minute
		if err != nil {
			log.Println(err)
		}

		go handleTest(conn)
	}
}

func handleTest(conn net.Conn) {
	defer conn.Close()
	httpHeaders := getHttpHeaders(conn)

	template, headerStatus := serveDirectory(httpHeaders)
	conn.Write([]byte(headerStatus + template + "\n\r"))
}

// Move all below this to a 'package parsers'
func getHttpHeaders(conn net.Conn) HttpHeaders {
	reader := bufio.NewReader(conn)
	var results []string
	moreLinesToRead := true
	for moreLinesToRead {
		line, _, err := reader.ReadLine()
		moreLinesToRead = strings.IndexByte(string(line), ' ') != -1
		if err != nil {
			log.Fatal("Couldn't read from connection")
		}

		results = append(results, string(line))
	}

	toReturn := HttpHeaders{}
	err := toReturn.fromStringArray(results)
	if err != nil {
		log.Fatal(err)
	}

	return toReturn
}

// TODO: make a config struct to make configuration more scalable
func loadConfig() Config {
	data, err := ioutil.ReadFile(CONFIG_PATH)
	if err != nil {
		log.Fatalf("Couldn't load config: %s", err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Couldn't load config: %s", err)
	}

	return config
}
