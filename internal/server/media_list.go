package server

import (
	"github.com/GaruGaru/ciak/internal/media/models"
	"github.com/GaruGaru/duty/task"
	"github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

type PageMedia struct {
	Media          models.Media
	TransferStatus task.ScheduledTask
	Cover          string
	Playable       bool
}

type MediaListPage struct {
	Title           string
	MediaCount      int
	PageMedia       []PageMedia
	NoMediasFound   bool
	TransferEnabled bool
}

func (p PageMedia) TButtonClass() string {
	switch p.TransferStatus.Status.State {
	case task.StateSuccess:
		return "btn-success"
	case task.StateError:
		return "btn-danger"
	case task.StateRunning:
		return "btn-secondary"
	case task.StatePending:
		return "btn-warning"
	default:
		return "btn-primary"
	}
}

func mediaToTitlesList(media []models.Media) []string {
	titles := make([]string, 0)
	for _, m := range media {
		titles = append(titles, m.Name)
	}
	return titles
}

func (s CiakServer) MediaListHandler(w http.ResponseWriter, r *http.Request) {

	mediaList, err := s.MediaDiscovery.Discover()

	mediaMetadata, err := s.OmbdClient.ByTitleBulk(mediaToTitlesList(mediaList)...)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	pageMediaList := make([]PageMedia, 0)

	for _, media := range mediaList {

		metadata := mediaMetadata[media.Name]

		if metadata.Poster == "" {
			metadata.Poster = "https://via.placeholder.com/300"
		}

		transferResult, _ := s.Daemon.Task(media.Hash())

		pageMediaList = append(pageMediaList, PageMedia{
			Media:          media,
			TransferStatus: transferResult,
			Cover:          metadata.Poster,
			Playable:       isExtensionPlayable(media.Format),
		})

	}

	mediaListPage := MediaListPage{
		Title:           "Home",
		MediaCount:      len(pageMediaList),
		PageMedia:       pageMediaList,
		NoMediasFound:   len(pageMediaList) == 0,
		TransferEnabled: s.Daemon.Conf.TransferDestination != "",
	}

	var mediaListTemplate = template.Must(template.ParseFiles("static/base.html", "static/media-list.html"))

	if err := mediaListTemplate.Execute(w, mediaListPage); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
