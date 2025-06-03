package ai_model

import (
	"bytes"
	"os"
	"text/template"
)

func LoadPrompt(filePath string, data any) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("prompt").Parse(string(content))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
