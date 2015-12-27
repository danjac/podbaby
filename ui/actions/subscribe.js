import * as api from '../api';
import { Actions } from '../constants';
import * as alerts from './alerts';
import { createAction } from './utils';

export function subscribe(id, name) {
  return dispatch => {
    api.subscribe(id);
    dispatch(createAction(Actions.SUBSCRIBE, id));
    dispatch(alerts.success(`You are now subscribed to ${name}`));
  };
}

export function unsubscribe(id, name) {
  return dispatch => {
    api.unsubscribe(id);
    dispatch(createAction(Actions.UNSUBSCRIBE, id));
    dispatch(alerts.success(`You are no longer subscribed to ${name}`));
  }
}
