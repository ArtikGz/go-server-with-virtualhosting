package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

// TODO: FIX DIRECTORY PATH TRAVERSAL
func serveDirectory(httpHeaders HttpHeaders) (string, string) {
	templateFolder, templateError := getTemplateFolderFromVirtualhost(httpHeaders.Host)
	templateFile := httpHeaders.Path
	if templateFile == "/" {
		templateFile = "/index.html" // TODO: This is temporary
	}
	template, err := readFile(templateFolder + templateFile)

	statusCode := 200
	if err != nil {
		template, err = readFile(templateFolder + templateError)
		statusCode = 404
		if err != nil {
			log.Fatal("You need to have an ErrorFile defined for each virtualhost on your config")
		}
	}

	var header string
	switch statusCode {
	case 200:
		header = HTTP_HEADER_STATUS_OK
	case 404:
		header = HTTP_HEADER_STATUS_NOT_FOUND
	default:
		header = HTTP_HEADER_STATUS_NOT_FOUND
	}

	return template, header
}

func getTemplateFolderFromVirtualhost(host string) (string, string) {
	for _, vhost := range GlobalConfig.VHosts {
		if vhost.Name == host {
			return vhost.Path, vhost.Default
		}
	}

	return GlobalConfig.Default, "/index.html"
}

func readFile(Path string) (string, error) {
	html, err := ioutil.ReadFile(Path)
	if err != nil {
		log.Printf("Can't read file: %s", err)
		return "Couldn't serve this to you :(", fmt.Errorf("file don't exist")
	}

	return string(html), nil
}
