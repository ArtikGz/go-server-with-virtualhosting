package main

import (
	"fmt"
	"strings"
)

type Address struct {
	host string
	port string
}

// TODO: Address.parseString()

type Match struct {
	Host string `json:"host"`
	Path string `json:"path"`
}

func (m Match) Matches(host string) bool {
	if m.Host == host {
		return true
	}

	return false
}

type HttpHeaders struct {
	Path string // GET /(path) HTTP/v
	Host string
}

func (httpHeaders *HttpHeaders) fromStringArray(headers []string) error {
	splitted := strings.Split(headers[0], " ") // headers[0] expected to be  'GET /(path) HTTP'
	if len(splitted) != 3 && splitted[0] == "GET" {
		return fmt.Errorf("Invalid header")
	}
	httpHeaders.Path = splitted[1]

	for _, header := range headers {
		if strings.HasPrefix(header, "Host") {
			splitted = strings.Split(header, " ")
			httpHeaders.Host = strings.Split(splitted[1], ":")[0]

		}
	}
	return nil
}

type Config struct {
	ErrorTemplate         string  `json:"ErrorTemplate"`
	DefaultTemplateFolder string  `json:"DefaultTemplateFolder"`
	Matchers              []Match `json:"Matchers"`
}
