import _ from 'lodash';

import { createAction } from './utils';

import { Actions } from '../constants';

export function addAlert(status, message) {
  return createAction(Actions.ADD_ALERT, {
    message,
    status,
    id: _.uniqueId()
  });
}

export const dismissAlert = id => createAction(Actions.DISMISS_ALERT, id);
