package main

import (
	"fmt"
	"github.com/gernest/front"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// recursively walk directory and return all files with given extension
func walk(root, ext string, index bool) (res []Link, i ContentIndex) {
	fmt.Printf("Scraping %s\n", root)
	i = make(ContentIndex)

	m := front.NewMatter()
	m.Handle("---", front.YAMLHandler)
	nPrivate := 0

	err := filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			res = append(res, parse(s, root)...)
			if index {
				text := getText(s)

				frontmatter, body, err := m.Parse(strings.NewReader(text))
				if err != nil {
					frontmatter = map[string]interface{}{}
					body = text
				}

				var title string
				if parsedTitle, ok := frontmatter["title"]; ok {
					title = parsedTitle.(string)
				} else {
					title = "Untitled Page"
				}

				// check if page is private
				if parsedPrivate, ok := frontmatter["draft"]; !ok || !parsedPrivate.(bool) {
					adjustedPath := strings.Replace(hugoPathTrim(trim(s, root, ".md")), " ", "-", -1)
					i[adjustedPath] = Content{
						Title:   title,
						Content: body,
					}
				} else {
					nPrivate++
				}
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Ignored %d private files \n", nPrivate)
	fmt.Printf("Parsed %d total links \n", len(res))
	return res, i
}

func getText(dir string) string {
	// read file
	fileBytes, err := ioutil.ReadFile(dir)
	if err != nil {
		panic(err)
	}

	return string(fileBytes)
}

