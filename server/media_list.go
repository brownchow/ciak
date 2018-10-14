package server

import (
	"github.com/GaruGaru/ciak/discovery"
	"html/template"
	"net/http"
)

type MediaListPage struct {
	MediaCount int
	MediaList  []discovery.Media
}

func (s CiakServer) MediaListHandler(w http.ResponseWriter, r *http.Request) {
	mediaList, err := s.MediaDiscovery.Discover()

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	pageTemplate, err := template.ParseFiles("static/media-list.html")

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	page := MediaListPage{
		MediaCount: len(mediaList),
		MediaList:  mediaList,
	}

	pageTemplate.Execute(w, page)

}