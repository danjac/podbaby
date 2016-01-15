import _ from 'lodash';

import { createAction } from './utils';
import { Actions, Alerts } from '../constants';

export function addAlert(status, message) {
  return createAction(Actions.ADD_ALERT, {
    message,
    status,
    id: _.uniqueId(),
  });
}

export const success = _.partial(addAlert, Alerts.SUCCESS);
export const info = _.partial(addAlert, Alerts.INFO);
export const warning = _.partial(addAlert, Alerts.WARNING);
export const danger = _.partial(addAlert, Alerts.DANGER);

export const dismissAlert = id => createAction(Actions.DISMISS_ALERT, id);
