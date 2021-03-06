package server

import (
	"fmt"
	"net/http"

	"github.com/GaruGaru/ciak/pkg/server/auth"
	"github.com/GaruGaru/ciak/pkg/utils"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

var UnauthenticatedUrls = []string{
	"/login",
	"/probe",
	"/api/login",
}

type LoginPage struct {
	Title string
}

var store = sessions.NewCookieStore([]byte("ciak_session"))

func (s CiakServer) LoginApiHandler(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")

	authUser, err := s.Authenticator.Authenticate(username, password)

	if err == nil {
		if err := s.createSession(w, r, authUser); err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

}

func (s CiakServer) createSession(w http.ResponseWriter, r *http.Request, user auth.User) error {
	session, err := store.Get(r, "user")

	if err != nil {
		return err
	}

	session.Values["username"] = user.Name
	return store.Save(r, w, session)
}

func (s CiakServer) LoginPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := loginTemplate.Execute(w, LoginPage{
			Title: "Login",
		})

		if err != nil {
			logrus.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s CiakServer) SessionAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !s.Config.AuthenticationEnabled || utils.StringIn(r.URL.Path, UnauthenticatedUrls) {
			next.ServeHTTP(w, r)
			return
		}

		session, err := store.Get(r, "user")

		if err != nil {
			fmt.Println("Session error ", err)
			return
		}

		if !session.IsNew {
			next.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
		}

	})
}
