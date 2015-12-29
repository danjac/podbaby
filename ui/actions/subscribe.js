import * as api from '../api';
import { Actions } from '../constants';
import * as alerts from './alerts';
import { createAction } from './utils';

export function subscribe(id) {
  return dispatch => {
    api.subscribe(id);
    dispatch(createAction(Actions.SUBSCRIBE, id));
    dispatch(alerts.success("You are now subscribed to this channel"));
  };
}

export function unsubscribe(id) {
  return dispatch => {
    api.unsubscribe(id);
    dispatch(createAction(Actions.UNSUBSCRIBE, id));
    dispatch(alerts.success("You are no longer subscribed to this channel"));
  }
}
