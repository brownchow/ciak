package server

import (
	"html/template"
	"net/http"

	"github.com/GaruGaru/ciak/pkg/config"
	"github.com/GaruGaru/ciak/pkg/media/details"
	"github.com/GaruGaru/ciak/pkg/media/discovery"
	"github.com/GaruGaru/ciak/pkg/server/auth"
	rice "github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var (
	loginTemplate     *template.Template
	mediaListTemplate *template.Template
)

const (
	serverVersion = "0.0.2"
)

type CiakServer struct {
	Config           config.CiakServerConfig
	MediaDiscovery   discovery.MediaDiscovery
	Authenticator    auth.Authenticator
	DetailsRetriever *details.Controller
}

func NewCiakServer(
	conf config.CiakServerConfig,
	discovery discovery.MediaDiscovery,
	authenticator auth.Authenticator,
	DetailsRetriever *details.Controller,
) CiakServer {
	server := CiakServer{
		Config:           conf,
		MediaDiscovery:   discovery,
		Authenticator:    authenticator,
		DetailsRetriever: DetailsRetriever,
	}
	// setup templates
	box := rice.MustFindBox("../../ui")

	loginTemplate := template.New("login")
	template.Must(loginTemplate.Parse(box.MustString("login.html")))
	template.Must(loginTemplate.Parse(box.MustString("base.html")))

	mediaListTemplate := template.New("help")
	template.Must(mediaListTemplate.Parse(box.MustString("media-list.html")))
	template.Must(mediaListTemplate.Parse(box.MustString("base.html")))

	return server
}

func (s CiakServer) Run() error {
	log.WithFields(log.Fields{
		"bind":    s.Config.ServerBinding,
		"version": serverVersion,
	}).Info("Ciak server started")

	router := mux.NewRouter()
	s.initRouting(router)
	return http.ListenAndServe(s.Config.ServerBinding, router)
}

func (s CiakServer) initRouting(router *mux.Router) {
	router.HandleFunc("/probe", ProbeHandler)
	router.HandleFunc("/", s.MediaListHandler())
	router.HandleFunc("/media/{media}", s.MediaStreamingHandler)
	router.HandleFunc("/login", s.LoginPageHandler())
	router.HandleFunc("/api/login", s.LoginApiHandler)
	router.Use(LoggingMiddleware)
	router.Use(s.SessionAuthMiddleware)
}
