import { Actions } from '../constants';

import * as api from '../api';
import { createAction } from './utils';

export function getLatestPodcasts() {
  return dispatch => {
    api.getLatestPodcasts()
    .then(result => {
      dispatch(createAction(Actions.LATEST_PODCASTS_SUCCESS, result.data));
    })
    .catch(error => {
        dispatch(createAction(Actions.LATEST_PODCASTS_FAILURE, { error }));
    });
  };
}
