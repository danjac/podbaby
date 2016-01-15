import * as api from '../api';
import { Actions } from '../constants';
import { createAction } from './utils';
import { requestPodcasts } from './podcasts';

export function searchChannel(query, id) {
  return dispatch => {
    dispatch(requestPodcasts());
    dispatch(createAction(Actions.CHANNEL_SEARCH_REQUEST, query));
    api.searchChannel(query, id)
    .then(result => {
      dispatch(createAction(Actions.CHANNEL_SEARCH_SUCCESS, result.data));
    })
    .catch(error => {
      dispatch(createAction(Actions.CHANNEL_SEARCH_FAILURE, { error }));
    });
  };
}

export function getChannel(id, page = 1) {
  return dispatch => {
    dispatch(requestPodcasts());
    dispatch(createAction(Actions.GET_CHANNEL_REQUEST));
    api.getChannel(id, page)
    .then(result => {
      dispatch(createAction(Actions.GET_CHANNEL_SUCCESS, result.data));
    })
    .catch(error => {
      dispatch(createAction(Actions.GET_CHANNEL_FAILURE, { error }));
    });
  };
}
