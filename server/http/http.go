package http

import (
	"html/template"
	"net/http"
	"os"

	"gopkg.in/macaron.v1"

	"github.com/Moekr/gopkg/logs"
	"github.com/Moekr/sword/common/args"
)

func StartHTTPService() {
	m := macaron.New()
	if os.Getenv("SWORD_DEV") != "" {
		m.Use(macaron.Logger())
	} else {
		macaron.Env = macaron.PROD
	}
	m.Use(macaron.Recovery())
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Directory: "tmpl",
		Funcs: []template.FuncMap{
			map[string]interface{}{
				"dict": dict,
			},
		},
	}))
	m.Use(macaron.Static("assets", macaron.StaticOptions{
		Prefix:      "assets",
		SkipLogging: macaron.Env == macaron.PROD,
	}))

	m.Get("/", handleIndexPage)
	m.Get("/index.html", handleIndexPage)
	m.Get("/detail.html", handleDetailPage)
	m.Group("/api", func() {
		m.Get("/conf", checkToken, handleAPIConf)
		m.Post("/push", checkToken, handleAPIPush)
		m.Get("/query/abbr", handleAPIAbbrQuery)
		m.Get("/query/full", handleAPIFullQuery)
		m.Get("/stat", handleAPIStat)
	}, setContentType, setErrorCode)

	go func() {
		logs.Info("[HTTP] service served on %s", args.Args.BindAddress)
		if err := http.ListenAndServe(args.Args.BindAddress, m); err != nil {
			logs.Panic("[HTTP] service serve error: %s", err.Error())
		}
	}()
}

func dict(values ...interface{}) map[string]interface{} {
	if len(values)%2 != 0 {
		return nil
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil
		}
		dict[key] = values[i+1]
	}
	return dict
}
