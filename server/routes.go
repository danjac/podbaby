package server

import (
	"database/sql"
	"errors"
	"github.com/danjac/podbaby/decoders"
	"github.com/danjac/podbaby/feedparser"
	"github.com/danjac/podbaby/models"
	"github.com/justinas/nosurf"
	"net/http"
)

func (s *Server) indexPage(w http.ResponseWriter, r *http.Request) {
	user, _ := s.getUserFromCookie(r)
	csrfToken := nosurf.Token(r)
	ctx := map[string]interface{}{
		"staticURL": s.Config.StaticURL,
		"csrfToken": csrfToken,
		"user":      user,
	}
	s.Render.HTML(w, http.StatusOK, "index", ctx)
}

func (s *Server) getLatestPodcasts(w http.ResponseWriter, r *http.Request) {
	user, _ := getUser(r)
	podcasts, err := s.DB.Podcasts.SelectAll(user.ID)
	if err != nil {
		s.abort(w, r, err)
		return
	}
	s.Render.JSON(w, http.StatusOK, podcasts)
}

func (s *Server) getChannels(w http.ResponseWriter, r *http.Request) {
	user, _ := getUser(r)
	channels, err := s.DB.Channels.SelectSubscribed(user.ID)
	if err != nil {
		s.abort(w, r, err)
		return
	}
	s.Render.JSON(w, http.StatusOK, channels)
}

func (s *Server) addChannel(w http.ResponseWriter, r *http.Request) {

	decoder := &decoders.NewChannel{}

	if err := decoders.Decode(r, decoder); err != nil {
		s.abort(w, r, HTTPError{http.StatusBadRequest, err})
		return
	}

	user, _ := getUser(r)

	go func(url string, userID int64) {

		result, err := feedparser.Fetch(url)

		if err != nil {
			s.Log.Error(err)
			return
		}

		channel := &models.Channel{
			URL:         url,
			Title:       result.Channel.Title,
			Image:       result.Channel.Image.Url,
			Description: result.Channel.Description,
		}

		if err := s.DB.Channels.Create(channel); err != nil {
			s.Log.Error(err)
			return
		}

		if err := s.DB.Subscriptions.Create(channel.ID, userID); err != nil {
			s.Log.Error(err)
			return
		}

		for _, item := range result.Items {
			podcast := &models.Podcast{
				ChannelID:   channel.ID,
				Title:       item.Title,
				Description: item.Description,
			}
			if len(item.Enclosures) == 0 {
				s.Log.Debug("Item has no enclosures")
				continue
			}
			podcast.EnclosureURL = item.Enclosures[0].Url
			pubDate, _ := item.ParsedPubDate()
			podcast.PubDate = pubDate

			if err := s.DB.Podcasts.Create(podcast); err != nil {
				s.Log.Error(err)
				return
			}
		}

	}(decoder.URL, user.ID)

	s.Render.Text(w, http.StatusCreated, "New channel")
}
func (s *Server) signup(w http.ResponseWriter, r *http.Request) {

	decoder := &decoders.Signup{}

	if err := decoders.Decode(r, decoder); r != nil {
		s.abort(w, r, HTTPError{http.StatusBadRequest, err})
		return
	}

	if exists, _ := s.DB.Users.IsEmail(decoder.Email); exists {
		s.abort(w, r, HTTPError{http.StatusBadRequest, errors.New("Email taken")})
		return
	}

	if exists, _ := s.DB.Users.IsName(decoder.Name); exists {
		s.abort(w, r, HTTPError{http.StatusBadRequest, errors.New("Name taken")})
		return
	}

	// make new user

	user := &models.User{
		Name:  decoder.Name,
		Email: decoder.Email,
	}

	if err := user.SetPassword(decoder.Password); err != nil {
		s.abort(w, r, err)
		return
	}

	if err := s.DB.Users.Create(user); err != nil {
		s.abort(w, r, err)
		return
	}
	s.setAuthCookie(w, user.ID)
	// tbd: no need to return user!
	s.Render.JSON(w, http.StatusCreated, user)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	decoder := &decoders.Login{}
	if err := decoders.Decode(r, decoder); err != nil {
		s.abort(w, r, HTTPError{http.StatusBadRequest, err})
		return
	}

	user, err := s.DB.Users.GetByNameOrEmail(decoder.Identifier)
	if err != nil {

		if err == sql.ErrNoRows {
			s.abort(w, r, HTTPError{http.StatusBadRequest, errors.New("no user found")})
			return
		}
		s.abort(w, r, err)
		return
	}

	if !user.CheckPassword(decoder.Password) {
		s.abort(w, r, HTTPError{http.StatusBadRequest, errors.New("Invalid password")})
		return
	}
	// login user
	s.setAuthCookie(w, user.ID)

	// tbd: no need to return user!
	s.Render.JSON(w, http.StatusOK, user)

}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	s.setAuthCookie(w, 0)
	s.Render.Text(w, http.StatusOK, "Logged out")
}
