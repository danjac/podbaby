import { pushPath } from 'redux-simple-router';

import * as api from '../api';
import { Actions } from '../constants';
import { createAction } from './utils';
import { requestPodcasts } from './podcasts';

export function search (query) {

  if (!query) {
    return createAction(Actions.SEARCH_SUCCESS, {
      channels: [],
      podcasts: [],
      isLoading: false
    });
  }
  return (dispatch, getState) => {
    dispatch(requestPodcasts());
    dispatch(createAction(Actions.SEARCH, query));
    api.search(query)
    .then(result => {
      dispatch(createAction(Actions.SEARCH_SUCCESS, result.data));
    })
    .catch(error => {
        dispatch(createAction(Actions.SEARCH_FAILURE, { error }));
    });
  }

}
