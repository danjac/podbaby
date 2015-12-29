package api

import "net/http"

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
