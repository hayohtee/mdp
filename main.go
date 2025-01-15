package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

//go:embed "templates"
var templateFS embed.FS

func main() {
	// Parse command-line flags.
	filename := flag.String("file", "", "Markdown file to preview")
	flag.Parse()

	// Check if the user provide the input file. If they did not, show usage.
	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run reads the content of the provided Markdown, convert it into
// an HTML format and save it with the same name as the Markdown.
func run(filename string) error {
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData, err := parseContent(input, "index.tmpl")
	if err != nil {
		return err
	}

	outputName := fmt.Sprintf("%s.html", filepath.Base(filename))
	fmt.Println(outputName)

	return saveHTML(outputName, htmlData)
}

// parseContent parses the content of the markdown file, sanitize it and
// include embed the result html content into an empty html page.
func parseContent(input []byte, templateFile string) ([]byte, error) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.Typographer),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithUnsafe(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert(input, &buf); err != nil {
		return nil, err
	}

	output := bluemonday.UGCPolicy().SanitizeReader(&buf)

	tmpl, err := template.New("markdown").ParseFS(templateFS, fmt.Sprintf("templates/%s", templateFile))
	if err != nil {
		return nil, err
	}

	var htmlBody bytes.Buffer
	if err := tmpl.ExecuteTemplate(&htmlBody, "htmlBody", output.Bytes()); err != nil {
		return nil, err
	}

	return htmlBody.Bytes(), nil
}
