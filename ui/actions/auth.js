import { pushPath } from 'redux-simple-router';

import * as api from '../api';
import { Actions } from '../constants';
import { createAction } from './utils';
import * as alerts from './alerts';

export const loginRequired = redirectTo => createAction(Actions.LOGIN_REQUIRED, redirectTo);
export const setCurrentUser = user => createAction(Actions.CURRENT_USER, user);

export function logout() {
  return dispatch => {
    api.logout();
    dispatch(createAction(Actions.LOGOUT));
    dispatch(pushPath("/"));
    dispatch(alerts.info("Bye for now!"))
  };
}

export function openRecoverPasswordForm() {
  return createAction(Actions.OPEN_RECOVER_PASSWORD_FORM);
}

export function closeRecoverPasswordForm() {
  return createAction(Actions.CLOSE_RECOVER_PASSWORD_FORM);
}

export function recoverPassword(identifier) {
  return dispatch => {
    dispatch(createAction(Actions.CLOSE_RECOVER_PASSWORD_FORM));
    api.recoverPassword(identifier)
    .then(result => {
      dispatch(createAction(Actions.RECOVER_PASSWORD_SUCCESS));
      dispatch(alerts.success("Please check your email inbox to recover your password"));
    })
    .catch(error => {
      dispatch(createAction(Actions.RECOVER_PASSWORD_FAILURE, { error }));
    });
  }
}
export function login(identifier, password) {
  return (dispatch, getState) => {
    dispatch(createAction(Actions.LOGIN));
    api.login(identifier, password)
    .then(result => {
      // call to api...
      const { auth } = getState();
      const nextPath = auth.redirectTo || '/podcasts/new/';
      dispatch(createAction(Actions.LOGIN_SUCCESS, result.data));
      dispatch(pushPath(nextPath));
      dispatch(alerts.success(`Welcome back, ${result.data.name}`))
    })
    .catch(error => {
      dispatch(createAction(Actions.LOGIN_FAILURE, { error }));
    });
  }
}

export function signup(name, email, password) {
  return dispatch =>  {
    dispatch(createAction(Actions.SIGNUP));
    api.signup(name, email, password)
    .then(result => {
      dispatch(createAction(Actions.SIGNUP_SUCCESS, result.data));
      dispatch(pushPath('/podcasts/new/'));
    })
    .catch(error => {
      dispatch(createAction(Actions.SIGNUP_FAILURE, { error }));
    });
  };
}
