package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/ShyBearStudio/tbot-admindashboard/logger"
)

func parseTemplateFiles(filenames ...string) (t *template.Template) {
	var files []string
	t = template.New("layout")
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	t = template.Must(t.ParseFiles(files...))
	return
}

func generateHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	if len(filenames) == 0 {
		if data != nil {
			logger.Warningln("File names list is empty but data is not empty")
		}
		return

	}
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", data)
}
