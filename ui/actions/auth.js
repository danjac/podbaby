import { pushPath } from 'redux-simple-router';

import * as api from '../api';
import { Actions } from '../constants';
import { createAction } from './utils';

export const loginRequired = redirectTo => createAction(Actions.LOGIN_REQUIRED, redirectTo);

export function logout() {
  return dispatch => {
    api.logout();
    dispatch(createAction(Actions.LOGOUT));
    dispatch(pushPath("/"));
  };
}

export function getCurrentUser() {
  return dispatch => {
    api.getCurrentUser()
    .then(result => dispatch(createAction(Actions.CURRENT_USER_SUCCESS, result.data)))
    .catch(()  => dispatch(createAction(Actions.CURRENT_USER_FAILURE)));
  };
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
    })
    .catch(error => {
      dispatch(createAction(Actions.LOGIN_FAILURE));
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
      dispatch(createAction(Actions.SIGNUP_FAILURE));
    });
  };
}
