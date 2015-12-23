import { pushPath } from 'redux-simple-router';

import * as api from '../api';
import { Actions } from '../constants';
import { createAction } from './utils';

export const loginRequired = redirectTo => createAction(Actions.LOGIN_REQUIRED, redirectTo);

export function logout() {
  return dispatch => {
    dispatch(createAction(Actions.LOGOUT));
    dispatch(pushPath("/"));
  };
}

export function login(identifier, password) {
  return (dispatch, getState) => {
    // call to api...
    const { auth } = getState();
    const nextPath = auth.redirectTo || '/podcasts/new/';

    dispatch(createAction(Actions.LOGIN_SUCCESS, { name: "danjac" }));
    dispatch(pushPath(nextPath));
  };
}

export function signup(name, email, password) {
  return dispatch =>  {
    dispatch(createAction(Actions.SIGNUP))
    api.signup(name, email, password)
    .then(result => dispatch(createAction(Actions.SIGNUP_SUCCESS, result.data)))
    .catch(dispatch(createAction(Actions.SIGNUP_FAILURE)));
  }
}
