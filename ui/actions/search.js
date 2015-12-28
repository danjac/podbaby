import { pushPath } from 'redux-simple-router';

import * as api from '../api';
import { Actions } from '../constants';
import { createAction } from './utils';

export function search (query) {

  return (dispatch, getState) => {
    dispatch(createAction(Actions.SEARCH, query));
    api.search(query)
    .then(result => {
      dispatch(createAction(Actions.SEARCH_SUCCESS, result.data));
      const { routing } = getState();
      if (routing.path !== "/podcasts/search/") {
        dispatch(pushPath("/podcasts/search/"));
      }
    })
    .catch(error => {
        dispatch(createAction(Actions.SEARCH_FAILURE, { error }));
    });
  }

}
