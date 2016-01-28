import * as api from '../api';
import { Actions } from '../constants';
import { createAction } from './utils';

export const requestPodcasts = () => createAction(Actions.PODCASTS_REQUEST);

export function getPodcast(id) {
  return (dispatch, getState) => {
    const { podcastCache } = getState().podcasts;
    const cached = podcastCache[id];
    if (cached) {
      dispatch(createAction(Actions.GET_PODCAST_SUCCESS, cached));
      return;
    }

    dispatch(createAction(Actions.GET_PODCAST_REQUEST));
    api.getPodcast(id)
    .then(result => {
      dispatch(createAction(Actions.GET_PODCAST_SUCCESS, result.data));
    })
    .catch(error => {
      dispatch(createAction(Actions.GET_PODCAST_FAILURE, { error }));
    });
  };
}
