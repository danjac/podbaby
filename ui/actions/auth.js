import { pushPath } from 'redux-simple-router';

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
    const nextPath = auth.redirectTo || '/';

    dispatch(createAction(Actions.LOGIN_SUCCESS, { username: "danjac" }));
    dispatch(pushPath(nextPath));
  };
}
