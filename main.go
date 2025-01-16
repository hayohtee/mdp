package main

import (
	"bytes"
	"embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"text/template"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
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

	if err := run(*filename, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run reads the content of the provided Markdown, convert it into
// an HTML format and save it in a temp folder and print the url
// to the generated html file to the stdout.
func run(filename string, out io.Writer) error {
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData, err := parseContent(input, "index.tmpl")
	if err != nil {
		return err
	}

	// Create a temp file using mdp prefix and .html suffix.
	temp, err := os.CreateTemp("", "mdp*.html")
	if err != nil {
		return err
	}

	// Close the temp file since we're not writing to it at the moment.
	if err := temp.Close(); err != nil {
		return err
	}

	outputName := temp.Name()
	fmt.Fprintln(out, outputName)

	return saveHTML(outputName, htmlData)
}

// parseContent parses the content of the markdown file, sanitize it and
// include embed the result html content into an empty html page.
func parseContent(input []byte, templateFile string) ([]byte, error) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.Typographer),
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
	if err := tmpl.ExecuteTemplate(&htmlBody, "htmlBody", output.String()); err != nil {
		return nil, err
	}

	return htmlBody.Bytes(), nil
}

// saveHTML saves the provided data to a file based on the provided filename.
func saveHTML(outputName string, data []byte) error {
	return os.WriteFile(outputName, data, 0644)
}

// preview open the provided filename using the default program.
func preview(filename string) error {
	var commandName string
	var commandParams []string

	switch runtime.GOOS {
	case "linux":
		commandName = "xdg-open"
	case "windows":
		commandName = "cmd.exe"
		commandParams = []string{"/C", "start"}
	case "darwin":
		commandName = "open"
	default:
		return errors.New("os not supported")
	}

	// Append the filename to command params
	commandParams = append(commandParams, filename)

	// Locate the executable in PATH
	cmdPath, err := exec.LookPath(commandName)
	if err != nil {
		return err
	}

	// Open the file using default program
	return exec.Command(cmdPath, commandParams...).Run()
}
