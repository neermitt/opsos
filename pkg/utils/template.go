package utils

import (
	"bytes"
	"text/template"
)

func ProcessTemplate(s string, vars map[string]any) (string, error) {
	tmpl, err := template.New("template").Parse(s)
	if err != nil {
		return "", err
	}
	var buff bytes.Buffer
	if err := tmpl.Execute(&buff, vars); err != nil {
		return "", err
	}
	return buff.String(), nil
}
