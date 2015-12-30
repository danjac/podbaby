import * as api from '../api';
import { Actions } from '../constants';
import { createAction } from './utils';
import { requestPodcasts } from './podcasts';

export function getChannel(id, page=1) {
    return dispatch => {
      dispatch(requestPodcasts());
      api.getChannel(id, page)
      .then(result => {
          dispatch(createAction(Actions.GET_CHANNEL_SUCCESS, result.data));
      })
      .catch(error => {
          dispatch(createAction(Actions.GET_CHANNEL_FAILURE, { error }));
      });
    };
}
