import * as api from '../api';
import { Actions } from '../constants';
import { createAction } from './utils';

export function getCategory(id) {
  return dispatch => {
    dispatch(createAction(Actions.GET_CATEGORY, id));
    dispatch(createAction(Actions.GET_CHANNELS_REQUEST));
    api.getCategory(id)
    .then(result => dispatch(createAction(Actions.GET_CHANNELS_SUCCESS, result.data)))
    .catch(error => dispatch(createAction(Actions.GET_CHANNELS_SUCCESS, { error })));
  };
}

export function loadCategories(categories) {
  return createAction(Actions.CATEGORIES_LOADED, categories);
}
