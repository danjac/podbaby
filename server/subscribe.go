package server

import "net/http"

func (s *Server) subscribe(w http.ResponseWriter, r *http.Request) {
	user, _ := getUser(r)
	channelID, _ := getInt64(r, "id")

	if err := s.DB.Subscriptions.Create(channelID, user.ID); err != nil {
		s.abort(w, r, err)
		return
	}
	s.Render.Text(w, http.StatusOK, "subscribed")
}

func (s *Server) unsubscribe(w http.ResponseWriter, r *http.Request) {
	user, _ := getUser(r)
	channelID, _ := getInt64(r, "id")

	if err := s.DB.Subscriptions.Delete(channelID, user.ID); err != nil {
		s.abort(w, r, err)
		return
	}
	s.Render.Text(w, http.StatusOK, "unsubscribed")
}
