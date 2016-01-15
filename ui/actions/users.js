import { pushPath } from 'redux-simple-router';

import * as api from '../api';
import { Actions } from '../constants';
import * as alerts from './alerts';
import { createAction } from './utils';


export function deleteAccount() {
  return dispatch => {
    api.deleteAccount()
    .then(() => {
      dispatch(createAction(Actions.DELETE_ACCOUNT_SUCCESS));
      dispatch(alerts.info('Your account has been deleted'));
      dispatch(pushPath('/'));
    })
    .catch(error => {
      dispatch(createAction(Actions.DELETE_ACCOUNT_FAILURE, { error }));
    });
  };
}

export function changeEmailComplete(email) {
  return dispatch => {
    dispatch(createAction(Actions.CHANGE_EMAIL_SUCCESS, email));
    dispatch(alerts.success('Your email has been updated'));
  };
}

export function changePasswordComplete() {
  return dispatch => {
    dispatch(createAction(Actions.CHANGE_PASSWORD_SUCCESS));
    dispatch(alerts.success('Your password has been updated'));
  };
}
