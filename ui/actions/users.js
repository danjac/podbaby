import * as api from '../api';
import { Actions } from '../constants';
import * as alerts from './alerts';
import { createAction } from './utils';


export function changeEmail(email) {
  return dispatch => {
    api.changeEmail(email)
    .then(result => {
      dispatch(createAction(Actions.CHANGE_EMAIL_SUCCESS, email));
      dispatch(alerts.success('Your email has been updated'));
    })
    .catch(error => {
      dispatch(createAction(Actions.CHANGE_EMAIL_FAILURE, { error }));
    });
  };
}

export function changePassword(oldPassword, newPassword) {
  return dispatch => {
    api.changePassword(oldPassword, newPassword)
      .then(result => {
        dispatch(createAction(Actions.CHANGE_PASSWORD_SUCCESS));
        dispatch(alerts.success('Your password has been updated'));
      })
      .catch(error => {
        dispatch(createAction(Actions.CHANGE_PASSWORD_FAILURE, { error }));
      });
  }
}
