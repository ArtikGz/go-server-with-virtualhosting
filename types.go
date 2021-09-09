package main

import (
	"fmt"
	"strings"
)

// TODO add TLS and http/2
type HttpHeaders struct {
	Path string // GET /(path) HTTP/v
	Host string
}

func (httpHeaders *HttpHeaders) fromStringArray(headers []string) error {
	splitted := strings.Split(headers[0], " ") // headers[0] expected to be  'GET /(path) HTTP'
	if len(splitted) != 3 && splitted[0] == "GET" {
		return fmt.Errorf("invalid header")
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
