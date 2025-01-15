package main

import (
	"flag"
	"fmt"
	"os"
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

	if err := run(*filename); err := nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}