package main

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var mds = `# header

Sample text.

[link](http://example.com)
`

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}
func WriteToFile(filepath string, data []byte) {

}

func FileExists(filepath string) bool {

	fileinfo, err := os.Stat(filepath)

	if os.IsNotExist(err) {
		return false
	}
	// Return false if the fileinfo says the file path is a directory.
	return !fileinfo.IsDir()
}

func mdToHtmlFull(filepath string) {
	// Check that the file is a markdown file
	if !FileExists(filepath) {
		log.Println("File path does not exists, aborting...")
		return
	}
	if path.Ext(filepath) != ".md" {
		log.Println("File is not markdown, ignoring...")
		return
	}

	// Read File into bytes
	b, err := os.ReadFile(filepath)
	if err != nil {
		log.Println("Error reading file, aborting. Error: \"", err, "\"")
	}

	// Convert to HTML
	html := mdToHTML(b)

	filename := strings.TrimSuffix(path.Base(filepath), path.Ext(filepath))
	target_filepath := "./html/" + filename + ".html"
	// Save to html
	os.WriteFile(target_filepath, html, 0777)
	log.Println("Converted", path.Base(filepath), "to html.")
}

func main() {
	input_directory := "./md"
	entries, err := os.ReadDir(input_directory)
	for _, entry := range entries {
		if !entry.IsDir() {
			mdToHtmlFull(path.Join(input_directory, entry.Name()))
		}
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					mdToHtmlFull(event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(input_directory)
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}
