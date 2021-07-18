package main

// TODO: Split the code into different files
// TODO: Move configuration to an specific struct

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

const (
	HTTP_HEADER_STATUS_OK = "HTTP/1.1 200 OK\n\r\n\r"
)

type Address struct {
	host string
	port string
}

type Match struct {
	Host     string `json:"host"`
	Filepath string `json:"filepath"`
}

func (m Match) Matches(addr Address) bool {
	host := addr.host
	if m.Host == host {
		return true
	}

	return false
}

func main() {
	result := loadConfig()
	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn, result)
	}
}

func handleConnection(conn net.Conn, matchers []Match) {
	log.Println(".handleConnection(): BEGIN")
	defer conn.Close()

	addr, err := getHostFromHTTP(conn)
	if err != nil {
		return
	}
	log.Printf(".handleConnection(): host -> %s", addr.host)

	msg := "<h1>Nothing</h1>"
	for _, matcher := range matchers {
		log.Println(matcher.Filepath)
		if matcher.Matches(addr) {
			msg = readHTML(matcher.Filepath)
			break
		}
	}

	conn.Write([]byte(HTTP_HEADER_STATUS_OK + msg + "\n\r"))
	log.Println(".handleConnection(): END")
}

func getHostFromHTTP(conn net.Conn) (Address, error) {
	log.Println(".getHostFromHTTP(): BEGIN")
	reader := bufio.NewReader(conn)

	hasHost := false
	for !hasHost {
		line, _, err := reader.ReadLine()
		log.Printf(".getHostFromHTTP(): READED LINE %s", string(line))
		if err != nil {
			return Address{}, fmt.Errorf("No line")
		}

		splitted := strings.Split(string(line), ":")
		if splitted[0] == "Host" {
			host := strings.Trim(splitted[1], " ")
			var port string
			if len(splitted) == 3 {
				port = strings.Trim(splitted[2], " ")
			}
			hasHost = true
			log.Println(".getHostFromHTTP(): END SUCCESSFUL")
			return Address{host, port}, nil
		}
	}

	log.Println(".getHostFromHTTP(): END WITH ERROR")
	return Address{"", ""}, fmt.Errorf("Couldn't resolve hostname")
}

func loadConfig() []Match {
	path := "./config.json"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Couldn't load config: %s", err)
	}
	log.Printf(".loadConfig(): RAW -> %s", data)

	var m []Match
	err = json.Unmarshal(data, &m)
	log.Printf(".loadConfig(): %s", m)
	if err != nil {
		log.Fatalf("Couldn't load config: %s", err)
	}

	return m
}

// TODO: Change to a readDirectory
// for serving the entire directory
func readHTML(Path string) string {
	html, err := ioutil.ReadFile(Path)
	if err != nil {
		log.Printf("Can't read file: %s", err)
		return "Couldn't serve this to you :("
	}

	return string(html)
}
