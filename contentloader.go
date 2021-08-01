package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

// TODO: FIX DIRECTORY PATH TRAVERSAL
func serveDirectory(httpHeaders HttpHeaders) string {
	templateFolder, err := getTemplateFolderFromVirtualhost(httpHeaders.Host)
	if err != nil {
		log.Fatal(err)
	}
	templateFile := httpHeaders.Path
	templateDir := templateFolder + templateFile
	template, err := readHTML(templateDir)
	if err != nil {
		template, err = readHTML(ERROR_PAGE_404)
		if err != nil {
			log.Fatal("You should have 404.html on templates folder")
		}
	}

	return template
}

// RETURN: HEADER - CONTENT - ERROR
func getTemplateFolderFromVirtualhost(host string) (string, error) {
	for _, matcher := range global_config {
		if matcher.Matches(host) {
			return matcher.Path, nil
		}
	}

	return "", fmt.Errorf("Not found")
}

func readHTML(Path string) (string, error) {
	html, err := ioutil.ReadFile(Path)
	if err != nil {
		log.Printf("Can't read file: %s", err)
		return "Couldn't serve this to you :(", fmt.Errorf("File don't exist")
	}

	return string(html), nil
}
