package server

import (
	"github.com/Moekr/sword/common"
	"html/template"
	"net/http"
)

var (
	htmlTemplate = template.New("template")
)

func init() {
	htmlTemplate.Parse(HeaderTemplate)
	htmlTemplate.Parse(FooterTemplate)
	htmlTemplate.Parse(IndexTemplate)
	htmlTemplate.Parse(DetailTemplate)
}

func httpIndex(w http.ResponseWriter, r *http.Request) {
	htmlTemplate.ExecuteTemplate(w, "index", conf.Targets)
}

func httpDetail(w http.ResponseWriter, r *http.Request) {
	targetId, err := parseIntParam(r, "t", false, -1)
	if err != nil {
		http.Redirect(w, r, "./", http.StatusMovedPermanently)
		return
	}
	var target *common.Target
	for _, tar := range conf.Targets {
		if tar.Id == targetId {
			target = tar
			break
		}
	}
	if target == nil {
		http.Redirect(w, r, "./", http.StatusMovedPermanently)
		return
	}
	htmlTemplate.ExecuteTemplate(w, "detail", map[string]interface{}{"target": target, "observers": conf.Observers})
}
