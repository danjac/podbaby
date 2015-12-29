package api

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"github.com/danjac/podbaby/decoders"
	"github.com/danjac/podbaby/feedparser"
	"github.com/danjac/podbaby/models"
	"github.com/justinas/nosurf"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const passwordChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateRandomPassword(length int) string {
	b := make([]byte, length)
	numChars := len(passwordChars)
	for i := range b {
		b[i] = passwordChars[rand.Intn(numChars)]
	}
	return string(b)
}

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
	result, err := s.DB.Podcasts.SelectSubscribed(user.ID, getPage(r))
	if err != nil {
		s.abort(w, r, err)
		return
	}
	s.Render.JSON(w, http.StatusOK, result)
}

func (s *Server) getChannelDetail(w http.ResponseWriter, r *http.Request) {
	user, _ := getUser(r)
	channelID, err := getInt64(r, "id")
	if err != nil {
		s.abort(w, r, err)
		return
	}
	channel, err := s.DB.Channels.GetByID(channelID, user.ID)
	if err != nil {
		s.abort(w, r, err)
		return
	}
	detail := &models.ChannelDetail{
		Channel: channel,
	}
	podcasts, err := s.DB.Podcasts.SelectByChannelID(channelID, user.ID, getPage(r))
	if err != nil {
		s.abort(w, r, err)
		return
	}
	for _, pc := range podcasts.Podcasts {
		pc.Name = channel.Title
		pc.Image = channel.Image
		pc.ChannelID = channel.ID
		detail.Podcasts = append(detail.Podcasts, pc)
	}
	detail.Page = podcasts.Page
	s.Render.JSON(w, http.StatusOK, detail)
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

func (s *Server) search(w http.ResponseWriter, r *http.Request) {

	user, _ := getUser(r)
	query := strings.Trim(r.FormValue("q"), " ")

	result := &models.SearchResult{}

	if query != "" {
		var err error
		if result.Channels, err = s.DB.Channels.Search(query, user.ID); err != nil {
			s.abort(w, r, err)
			return
		}
		if result.Podcasts, err = s.DB.Podcasts.Search(query, user.ID); err != nil {
			s.abort(w, r, err)
			return
		}
	}

	s.Render.JSON(w, http.StatusOK, result)
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

func (s *Server) recoverPassword(w http.ResponseWriter, r *http.Request) {
	// generate a temp password
	decoder := &decoders.RecoverPassword{}

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

	tempPassword := generateRandomPassword(6)

	user.SetPassword(tempPassword)

	if err := s.DB.Users.UpdatePassword(user.Password, user.ID); err != nil {
		s.abort(w, r, err)
		return
	}

	// send email to user
	// TBD: add email config to server config
	go func(r *http.Request, user *models.User, tempPassword string) {

		msg := fmt.Sprintf(`Hi %s,
We've reset your password so you can sign back in again!

Here is your new temporary password:

%s

You can login here:

%s/#/login/

Change your password as soon as possible!

Thanks,

PodBaby
    `, user.Name, tempPassword, r.Host)

		s.Log.Info(msg)

		err := smtp.SendMail(
			"mail.localhost:25",
			nil, // auth
			"sender@podbaby.me",
			[]string{user.Email},
			[]byte(msg),
		)

		if err != nil {
			s.Log.Error(err)
		}

	}(r, user, tempPassword)

	s.Render.Text(w, http.StatusOK, "password sent")

}

func (s *Server) signup(w http.ResponseWriter, r *http.Request) {

	decoder := &decoders.Signup{}

	if err := decoders.Decode(r, decoder); r != nil {
		s.abort(w, r, HTTPError{http.StatusBadRequest, err})
		return
	}

	if exists, _ := s.DB.Users.IsEmail(decoder.Email, 0); exists {
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

func (s *Server) subscribe(w http.ResponseWriter, r *http.Request) {
	user, _ := getUser(r)
	channelID, err := getInt64(r, "id")
	if err != nil {
		s.abort(w, r, err)
		return
	}
	if err := s.DB.Subscriptions.Create(channelID, user.ID); err != nil {
		s.abort(w, r, err)
		return
	}
	s.Render.Text(w, http.StatusOK, "subscribed")
}

func (s *Server) unsubscribe(w http.ResponseWriter, r *http.Request) {
	user, _ := getUser(r)
	channelID, err := getInt64(r, "id")
	if err != nil {
		s.abort(w, r, err)
		return
	}
	if err := s.DB.Subscriptions.Delete(channelID, user.ID); err != nil {
		s.abort(w, r, err)
		return
	}
	s.Render.Text(w, http.StatusOK, "unsubscribed")
}

func (s *Server) getBookmarks(w http.ResponseWriter, r *http.Request) {
	user, _ := getUser(r)
	result, err := s.DB.Podcasts.SelectBookmarked(user.ID, getPage(r))
	if err != nil {
		s.abort(w, r, err)
		return
	}
	s.Render.JSON(w, http.StatusOK, result)
}

func (s *Server) addBookmark(w http.ResponseWriter, r *http.Request) {
	user, _ := getUser(r)
	podcastID, err := getInt64(r, "id")
	if err != nil {
		s.abort(w, r, err)
		return
	}
	if err := s.DB.Bookmarks.Create(podcastID, user.ID); err != nil {
		s.abort(w, r, err)
		return
	}
	s.Render.Text(w, http.StatusOK, "bookmarked")
}

func (s *Server) removeBookmark(w http.ResponseWriter, r *http.Request) {
	user, _ := getUser(r)
	podcastID, err := getInt64(r, "id")
	if err != nil {
		s.abort(w, r, err)
		return
	}
	if err := s.DB.Bookmarks.Delete(podcastID, user.ID); err != nil {
		s.abort(w, r, err)
		return
	}
	s.Render.Text(w, http.StatusOK, "bookmark removed")
}

func (s *Server) changeEmail(w http.ResponseWriter, r *http.Request) {
	user, _ := getUser(r)
	decoder := &decoders.NewEmail{}
	if err := decoders.Decode(r, decoder); err != nil {
		s.abort(w, r, HTTPError{http.StatusBadRequest, err})
		return
	}
	// does this email exist?
	if exists, _ := s.DB.Users.IsEmail(decoder.Email, user.ID); exists {
		s.abort(w, r, HTTPError{http.StatusBadRequest, errors.New("Email taken")})
		return
	}

	if err := s.DB.Users.UpdateEmail(decoder.Email, user.ID); err != nil {
		s.abort(w, r, err)
		return
	}
	s.Render.Text(w, http.StatusOK, "email updated")
}

func (s *Server) changePassword(w http.ResponseWriter, r *http.Request) {
	user, _ := getUser(r)
	decoder := &decoders.NewPassword{}
	if err := decoders.Decode(r, decoder); err != nil {
		s.abort(w, r, HTTPError{http.StatusBadRequest, err})
		return
	}

	// validate old password first

	if !user.CheckPassword(decoder.OldPassword) {
		s.abort(w, r, HTTPError{http.StatusBadRequest, errors.New("Invalid password")})
		return
	}
	user.SetPassword(decoder.NewPassword)

	if err := s.DB.Users.UpdatePassword(user.Password, user.ID); err != nil {
		s.abort(w, r, err)
		return
	}
	s.Render.Text(w, http.StatusOK, "password updated")
}
