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
		template, err = readFile(templateError)
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

// RETURN template, errorTemplate
func getTemplateFolderFromVirtualhost(host string) (string, string) {
	for _, matcher := range GlobalConfig.Matchers {
		if matcher.Matches(host) {
			return matcher.Path, matcher.ErrorFile
		}
	}

	defaultMatcher := GlobalConfig.DefaultMatcher
	return defaultMatcher.Path, defaultMatcher.ErrorFile
}

func readFile(Path string) (string, error) {
	html, err := ioutil.ReadFile(Path)
	if err != nil {
		log.Printf("Can't read file: %s", err)
		return "Couldn't serve this to you :(", fmt.Errorf("File don't exist")
	}

	return string(html), nil
}
