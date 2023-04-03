package web

import (
	"embed"
	"github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

//go:embed template
var templateFiles embed.FS

type templateData struct {
	Hostname string
}

func RootHandler() http.Handler {

	t, err := template.ParseFS(templateFiles, "template/home.html")

	if err != nil {
		logrus.WithError(err).Error("Could not parse template")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")

		data := templateData{
			Hostname: r.Host,
		}

		err = t.ExecuteTemplate(w, "home.html", data)
		if err != nil {
			logrus.WithError(err).Error("Could generate page")
		}
		logrus.Info("index")
	})
}
