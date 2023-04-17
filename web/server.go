package web

import (
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
	"poulpe.ztec.fr/engine"
)

func StartServer(addr string, searchEngine engine.Engine) error {
	router := httprouter.New()

	router.Handler("GET", "/", PrometheusInstrumentation("root", RootHandler()))
	router.Handler("GET", "/statics/*filepath", PrometheusInstrumentation("static", NewStaticHandler()))
	router.Handler("GET", "/search", PrometheusInstrumentation("search", NewInstantSearchHandler(searchEngine)))
	router.Handler("GET", "/metrics", PrometheusInstrumentation("metrics", promhttp.Handler()))

	logrus.Infof("Server listening on %s", addr)
	return http.ListenAndServe(addr, router)
}
