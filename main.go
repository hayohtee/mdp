package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)


func main()  {
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

func run(filename string) error {
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData := parseContent(input)

	outputName := fmt.Sprintf("%s.html", filepath.Base(filename))
	fmt.Println(outputName)

	return saveHTML(outputName, htmlData)
}