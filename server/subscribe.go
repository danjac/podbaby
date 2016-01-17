package server

import "net/http"

func subscribe(s *Server, w http.ResponseWriter, r *http.Request) error {
	user, _ := getUser(r)
	channelID, _ := getID(r)

	if err := s.DB.Subscriptions.Create(channelID, user.ID); err != nil {
		return err
	}
	return s.Render.Text(w, http.StatusOK, "subscribed")
}

func unsubscribe(s *Server, w http.ResponseWriter, r *http.Request) error {
	user, _ := getUser(r)
	channelID, _ := getID(r)

	if err := s.DB.Subscriptions.Delete(channelID, user.ID); err != nil {
		return err
	}
	return s.Render.Text(w, http.StatusOK, "unsubscribed")
}
