package processor

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"text/template"
)

type Processor struct {
	model Model
}

type Model struct {
	Var map[string]string
}

func New() Processor {
	return Processor{}
}

func NewWithModel(m Model) Processor {
	proc := Processor{}
	proc.model = m
	return proc
}

func (p *Processor) SetModel(m Model) {
	p.model = m
}

// ParseAndExecutePath will parse the path as a template and execute it using the settings provided
func (p Processor) ParseAndExecutePath(path string) (string, error) {
	mainTemplate := template.New("main")

	tmpl, err := mainTemplate.Parse(path)
	if err != nil {
		return "", fmt.Errorf("error parsing path '%v' to template: %w", path, err)
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, p.model)
	if err != nil {
		return "", fmt.Errorf("error executing path '%v' as template: %w", path, err)
	}

	return buf.String(), nil
}

// ParseAndExecuteFile will parse a file as a template and execute it using the Model set. it will write out to the writer once Executed.
func (p Processor) ParseAndExecuteFile(sourcePath string, wr io.Writer) error {
	fileTemplate, err := template.ParseFiles(sourcePath)
	if err != nil {
		return fmt.Errorf("error Parsing template for file '%v': %w", sourcePath, err)
	}

	if err = fileTemplate.ExecuteTemplate(wr, filepath.Base(sourcePath), p.model); err != nil {
		return fmt.Errorf("error executing template file '%v': %w", sourcePath, err)
	}

	return nil
}
