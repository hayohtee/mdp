# mdp
mdp is a command-line tool that converts a Markdown file into HTML format that can be viewed in a browser.

## Features
- Convert Markdown file into HTML format
- Sanitize the generated HTML document
- Previewing the generated document in the browser.

## Building from source
1. Clone the repository
   ```bash
   git clone git@github.com:hayohtee/mdp.git
   ```
2. Change into the project directory
   ```bash
   cd mdp
   ```
3. Compile
   ```bash
   go build ./...
   ```

## Usage
1. Convert a Markdown\
   To convert a and preview a markdown simply use the -file flag and specify the Markdown file that you want to convert\
   Here is an example showing how to convert README.md file
   ```bash
   ./mdp -file README.md
   ```
2. Disable auto preview\
   By default mdp auto preview the generated HTML file in user's default program associated with HTML file.\
   To disable this feature, include -s in addition to -file flag which skip autopreview.\
   Here is an example showing how to convert README.md file while skipping auto preview
   ```bash
   ./mdp -file README.md -s
   ```
3. Show all available options
   ```bash
   ./todo -h
   ```
