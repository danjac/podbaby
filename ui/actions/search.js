import * as api from '../api';
import { Actions } from '../constants';
import { createAction } from './utils';
import { requestPodcasts } from './podcasts';

export function clearSearch() {
  return createAction(Actions.CLEAR_SEARCH);
}

export function changeSearchType(type) {
  return createAction(Actions.CHANGE_SEARCH_TYPE, type);
}

export function search(query, type, page = 1) {
  if (!query) {
    return clearSearch();
  }

  return dispatch => {
    dispatch(requestPodcasts());
    dispatch(createAction(Actions.SEARCH_REQUEST, { query, type }));
    api.search(query, type, page)
    .then(result => {
      dispatch(createAction(Actions.SEARCH_SUCCESS, result.data));
    })
    .catch(error => {
      dispatch(createAction(Actions.SEARCH_FAILURE, { error }));
    });
  };
}
