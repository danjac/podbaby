import { Actions } from '../constants';
import * as api from '../api';
import * as alerts from './alerts';
import { requestPodcasts } from './podcasts';
import { createAction } from './utils';


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
  dispatch(createAction(Actions.PODCASTS_REQUEST));
  return dispatch => {
    dispatch(requestPodcasts());
    api.getBookmarks(page)
    .then(result => {
      dispatch(createAction(Actions.GET_BOOKMARKS_SUCCESS, result.data));
    })
    .catch(error => {
      dispatch(createAction(Actions.GET_BOOKMARKS_FAILURE, { error }));
    });
  }
}
