import * as api from '../api';

import { Actions } from '../constants';
import { createAction } from './utils';
import * as alerts from './alerts';

export function clearAll() {
  api.clearAllPlayed();
  return dispatch => {
    dispatch(alerts.success("All podcasts have been removed from your recently played list"));
    dispatch(createAction(Actions.CLEAR_RECENT_PLAYS));
  };
}

export function getRecentlyPlayed(page=1) {
  return dispatch => {
    dispatch(createAction(Actions.PODCASTS_REQUEST));
    api.getRecentlyPlayed(page)
    .then(result => {
      dispatch(createAction(Actions.GET_RECENT_PLAYS_SUCCESS, result.data));
    })
    .catch(error => {
      dispatch(createAction(Actions.GET_RECENT_PLAYS_FAILURE, { error }));
    });
  };
}
