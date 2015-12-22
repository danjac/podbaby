import { pushPath } from 'redux-simple-router';

import { Actions } from '../constants';
import { createAction } from './utils';

export function search (query) {
  return (dispatch, getState) => {
    dispatch(createAction(Actions.SEARCH, query));
    const { routing } = getState();
    if (routing.path !== "/podcasts/search/") {
      dispatch(pushPath("/podcasts/search/"));
    }
  }

}
