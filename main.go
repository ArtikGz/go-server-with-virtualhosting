package main

// TODO: Make templates folder configurable
// TODO: Serve index.html if no filename provided

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
	ERROR_PAGE_404               = "templates/404.html" // TODO: Make this configurable
)

var global_config []Match

func main() {
	global_config = loadConfig()
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

	template := serveDirectory(httpHeaders)
	conn.Write([]byte(HTTP_HEADER_STATUS_OK + template + "\n\r"))
}

// Move all below this to a 'package parsers'
// TODO: Move this to a parsers.go
func getHttpHeaders(conn net.Conn) HttpHeaders {
	log.Println(".getHttpHeaders(): BEGIN")
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
func loadConfig() []Match {
	data, err := ioutil.ReadFile(CONFIG_PATH)
	if err != nil {
		log.Fatalf("Couldn't load config: %s", err)
	}

	var hosts []Match
	err = json.Unmarshal(data, &hosts)
	log.Printf(".loadConfig(): %s", hosts)
	if err != nil {
		log.Fatalf("Couldn't load config: %s", err)
	}

	return hosts
}
