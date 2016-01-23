import * as api from '../api';
import { Actions } from '../constants';
import { createAction } from './utils';
import { requestPodcasts } from './podcasts';

export function clearSearch() {
  return createAction(Actions.CLEAR_SEARCH);
}

export function search(query) {
  if (!query) {
    return clearSearch();
  }

  return dispatch => {
    dispatch(requestPodcasts());
    dispatch(createAction(Actions.SEARCH_REQUEST, query));
    api.search(query)
    .then(result => {
      dispatch(createAction(Actions.SEARCH_SUCCESS, result.data));
    })
    .catch(error => {
      dispatch(createAction(Actions.SEARCH_FAILURE, { error }));
    });
  };
}
