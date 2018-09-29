package server

import (
	"html/template"
	"net/http"
)

var (
	indexTemplate = template.New("index")
)

func init() {
	indexTemplate.Parse(IndexTemplate)
}

func httpIndex(w http.ResponseWriter, r *http.Request) {
	indexTemplate.Execute(w, conf.Targets)
}
