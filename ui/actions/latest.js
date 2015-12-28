import { Actions } from '../constants';

import * as api from '../api';
import { createAction } from './utils';

export function getLatestPodcasts(page=1) {
  return dispatch => {
    api.getLatestPodcasts(page)
    .then(result => {
      dispatch(createAction(Actions.LATEST_PODCASTS_SUCCESS, result.data));
    })
    .catch(error => {
        dispatch(createAction(Actions.LATEST_PODCASTS_FAILURE, { error }));
    });
  };
}
