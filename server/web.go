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
	htmlTemplate.Parse(HeadTemplate)
	htmlTemplate.Parse(HeaderTemplate)
	htmlTemplate.Parse(FooterTemplate)
	htmlTemplate.Parse(IndexTemplate)
	htmlTemplate.Parse(DetailTemplate)
}

func httpIndex(w http.ResponseWriter, r *http.Request) {
	timeRange, err := parseIntParam(r, "r", true, 1)
	if err != nil {
		timeRange = rangeDay
	}
	params := map[string]interface{}{
		"targets":   conf.Targets,
		"timeRange": timeRange,
	}
	htmlTemplate.ExecuteTemplate(w, "index", params)
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
	timeRange, err := parseIntParam(r, "r", true, 1)
	if err != nil {
		timeRange = rangeDay
	}
	params := map[string]interface{}{
		"target":    target,
		"observers": conf.Observers,
		"timeRange": timeRange,
	}
	htmlTemplate.ExecuteTemplate(w, "detail", params)
}
