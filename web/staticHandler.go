package web

import (
	"embed"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"poulpe.ztec.fr/engine"
	"poulpe.ztec.fr/types"
)

//go:embed statics
var staticFiles embed.FS

func NewStaticHandler() http.Handler {
	var staticFS = http.FS(staticFiles)
	return http.FileServer(staticFS)
}

func NewInstantSearchHandler(searchEngine engine.Engine) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query()
		searchString := query.Get("q")
		var result []types.EmojiDescription

		result, err := searchEngine.Search(searchString)
		if err != nil {
			logrus.WithField("q", searchString).WithError(err).Error("could not search")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("ERROR"))
		} else {
			resultJson, err := json.Marshal(result)
			if err != nil {
				logrus.WithField("q", searchString).WithError(err).Error("could not build response")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("ERROR"))
				return
			}
			logrus.WithField("q", searchString).WithField("results", len(result)).Info("search")
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(resultJson)
		}
	})
}
