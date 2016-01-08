import { Actions } from '../constants';
import * as api from '../api';
import * as alerts from './alerts';
import { requestPodcasts } from './podcasts';
import { createAction } from './utils';

export function toggleBookmark(podcast) {
  return podcast.isBookmarked ? deleteBookmark(podcast.id) : addBookmark(podcast.id);
}

export function addBookmark(podcastId) {
  return dispatch => {
    api.addBookmark(podcastId);
    dispatch(createAction(Actions.ADD_BOOKMARK, podcastId));
    dispatch(alerts.success("You have bookmarked this podcast"));
  };
}

export function deleteBookmark(podcastId) {
  return dispatch => {
    api.deleteBookmark(podcastId);
    dispatch(createAction(Actions.DELETE_BOOKMARK, podcastId));
    dispatch(alerts.success("You have removed this bookmark"));
  };
}

export function getBookmarks(page=1) {
  return dispatch => {
    dispatch(requestPodcasts());
    dispatch(createAction(Actions.CLEAR_BOOKMARKS_SEARCH));
    api.getBookmarks(page)
    .then(result => {
      dispatch(createAction(Actions.GET_BOOKMARKS_SUCCESS, result.data));
    })
    .catch(error => {
      dispatch(createAction(Actions.GET_BOOKMARKS_FAILURE, { error }));
    });
  }
}

export function searchBookmarks(query) {
  return dispatch => {
      dispatch(requestPodcasts());
      dispatch(createAction(Actions.BOOKMARKS_SEARCH_REQUEST, query));
      api.searchBookmarks(query)
      .then(result => {
        dispatch(createAction(Actions.BOOKMARKS_SEARCH_SUCCESS, result.data));
      })
      .catch(error => {
          dispatch(createAction(Actions.BOOKMARKS_SEARCH_FAILURE, { error }));
      });
  };
}
