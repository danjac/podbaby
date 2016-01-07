import * as api from '../api';
import { Actions } from '../constants';
import { createAction } from './utils';

export const requestPodcasts = () => createAction(Actions.PODCASTS_REQUEST);

export function getPodcast (id) {
  return dispatch => {
    dispatch(createAction(Actions.GET_PODCAST_REQUEST));
    api.getPodcast(id)
    .then(result => {
      dispatch(createAction(Actions.GET_PODCAST_SUCCESS, result.data));
    })
    .catch(error => {
      dispatch(createAction(Actions.GET_PODCAST_FAILURE, { error }));
    });
  };
};
