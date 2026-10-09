package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"text/template"

	"charm.land/glamour/v2"
)

//go:embed help.tmpl
var helpTemplate string

const (
	glamourStyle string = "dracula"
	glamourWidth int    = 80
)

func printHelp(block string) error {

	tmpl, err := template.New("help").Parse(helpTemplate)
	if err != nil {
		return fmt.Errorf("failed to analyze template : %v", err)
	}

	// Vérifier si le bloc existe dans le template
	if tmpl.Lookup(block) == nil {
		block = "erreur_404"
	}

	// Évaluer le template en mémoire
	var rawBuffer bytes.Buffer
	err = tmpl.ExecuteTemplate(&rawBuffer, block, nil)
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
