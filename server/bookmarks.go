package server

import "net/http"

func getBookmarks(s *Server, w http.ResponseWriter, r *http.Request) error {
	user, _ := getUser(r)
	result, err := s.DB.Podcasts.SelectBookmarked(user.ID, getPage(r))
	if err != nil {
		return err
	}
	return s.Render.JSON(w, http.StatusOK, result)
}

func addBookmark(s *Server, w http.ResponseWriter, r *http.Request) error {
	user, _ := getUser(r)
	podcastID, _ := getID(r)

	if err := s.DB.Bookmarks.Create(podcastID, user.ID); err != nil {
		return err
	}
	return s.Render.Text(w, http.StatusOK, "bookmarked")
}

func removeBookmark(s *Server, w http.ResponseWriter, r *http.Request) error {
	user, _ := getUser(r)
	podcastID, _ := getID(r)

	if err := s.DB.Bookmarks.Delete(podcastID, user.ID); err != nil {
		return err
	}
	return s.Render.Text(w, http.StatusOK, "bookmark removed")
}
