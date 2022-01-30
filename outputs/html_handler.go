package outputs

import (
	"bytes"
	"embed"
	"html/template"
	"streamjury/gameplay"
)

const (
	filename = "round.html"
)

//go:embed round.html
var dataTemplateFS embed.FS

func PublishResultsInHTML(g gameplay.GamePlay) ([]byte, error) {
	var err error
	var tpl *template.Template
	var buf bytes.Buffer = bytes.Buffer{}

	tpl, err = template.ParseFS(dataTemplateFS, filename)
	if err != nil {
		return []byte{}, err
	}
	if err = tpl.ExecuteTemplate(&buf, filename, g); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}
