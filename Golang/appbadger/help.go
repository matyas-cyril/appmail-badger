package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"text/template"
	"time"

	"charm.land/glamour/v2"
)

//go:embed help.tmpl
var helpTemplate string

const (
	glamourStyle string = "dracula"
	glamourWidth int    = 150
)

func printHelp(block string) error {

	type Tmpl struct {
		Conf  string
		Epoch int64
	}

	tmpl, err := template.New("help").Parse(helpTemplate)
	if err != nil {
		return fmt.Errorf("failed to analyze template : %v", err)
	}

	// Vérifier si le bloc existe dans le template
	if tmpl.Lookup(block) == nil {
		block = "erreur_404"
	}

	data := Tmpl{
		Conf:  DEFAULT_CONFIG_FILE,
		Epoch: time.Now().Unix(),
	}

	// Évaluer le template en mémoire
	var rawBuffer bytes.Buffer
	err = tmpl.ExecuteTemplate(&rawBuffer, block, data)
	if err != nil {
		return fmt.Errorf("failed to execute template : %v", err)
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithWordWrap(glamourWidth),
		glamour.WithStandardStyle(glamourStyle),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize Glamour librairy : %v", err)
	}

	rendered, err := r.Render(rawBuffer.String())
	if err != nil {
		return fmt.Errorf("failed to Markdown rendering : %v", err)
	}

	fmt.Fprint(os.Stdout, rendered)

	return nil
}
