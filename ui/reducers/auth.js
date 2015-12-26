import _ from 'lodash';
import { Actions } from '../constants';

const initialState = {
  isLoggedIn: false,
  name: null,
  redirectTo: null,
  isLoaded: true
};

export default function(state=initialState, action) {
  switch(action.type) {
    case Actions.LOGIN_REQUIRED:
      return Object.assign({}, state, { redirectTo: action.payload });
    case Actions.LOGIN_SUCCESS:
    case Actions.SIGNUP_SUCCESS:
      return Object.assign({}, state, action.payload, { isLoggedIn: true });
    case Actions.CURRENT_USER_SUCCESS:
      return Object.assign({}, state, action.payload, { isLoggedIn: !_.isEmpty(action.payload) });

    case Actions.LOGOUT:
      return initialState;
  }
  return state;
}
