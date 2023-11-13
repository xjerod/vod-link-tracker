package main

import (
	"os"
	"path"
	"text/template"
	"time"
)

type Post struct {
	Name string
	Date time.Time
}

type Content struct {
	Name  string
	Date  time.Time
	Links []string
}

func main() {
	links := Content{
		Name: "UAW VICTORY!!! DAGESTAN ANTISEMITISM! GAZA WARCRIMES CONTINUE.",
		Date: time.Now(),
		Links: []string{
			"twitter.com/kaitlancollins/status/1717245490826272911",
		},
	}

	tmplFile := "./azanlinks/archetypes/links.md"
	name := path.Base(tmplFile)
	tmpl := template.Must(template.New(name).ParseFiles(tmplFile))

	file, err := os.Create("./azanlinks/content/posts/2023-10-29.md")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = tmpl.Execute(file, links)
	if err != nil {
		panic(err)
	}
}
