import * as api from '../api';
import { Actions } from '../constants';
import { createAction } from './utils';

const fetchChannels = apiCall => {
  return dispatch => {
    dispatch(createAction(Actions.GET_CHANNELS_REQUEST));
    apiCall()
    .then(result => {
      dispatch(createAction(Actions.GET_CHANNELS_SUCCESS, result.data));
    })
    .catch(error => {
      dispatch(createAction(Actions.GET_CHANNELS_FAILURE, { error }));
    });
  };
};

export function filterChannels(filter) {
  return createAction(Actions.FILTER_CHANNELS, filter);
}

export function selectPage(page) {
  return createAction(Actions.SELECT_CHANNELS_PAGE, page);
}

export function getRecommendations() {
  return fetchChannels(api.getRecommendations);
}

export function getSubscriptions() {
  return fetchChannels(api.getSubscriptions);
}
