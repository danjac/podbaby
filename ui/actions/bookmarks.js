import { Actions } from '../constants';
import * as api from '../api';
import * as alerts from './alerts';

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

export function getBookmarks() {
  return dispatch => {
    api.getBookmarks()
    .then(result => {
      dispatch(createAction(Actions.GET_BOOKMARKS_SUCCESS, result.data));
    })
    .catch(error => {
      dispatch(createAction(Actions.GET_BOOKMARKS_FAILURE, { error }));
    });
  }
}
