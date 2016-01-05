import * as api from '../api';

import { Actions } from '../constants';
import { createAction } from './utils';

export function getRecentlyPlayed() {
  return dispatch => {
    dispatch(createAction(Actions.PODCASTS_REQUEST));
    api.getRecentlyPlayed()
    .then(result => {
      dispatch(createAction(Actions.GET_RECENT_PLAYS_SUCCESS, result.data));
    })
    .catch(error => {
      dispatch(createAction(Actions.GET_RECENT_PLAYS_FAILURE, { error }));
    });
  };
}
