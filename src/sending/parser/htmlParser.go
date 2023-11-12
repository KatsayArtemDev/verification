package parser

import (
	"bytes"
	"fmt"
	"html/template"
)

type emailContent struct {
	Number int
	Pin    string
}

func HtmlParser(templatePath, pin string) (string, error) {
	var body bytes.Buffer
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("error parsing data from html file: %v", err)
	}

	err = t.Execute(&body, emailContent{1, pin})
	if err != nil {
		return "", fmt.Errorf("error executing data from html file: %v", err)
	}

	return body.String(), nil
}
