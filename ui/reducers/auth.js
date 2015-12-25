import { Actions } from '../constants';

const initialState = {
  isLoggedIn: false,
  name: null,
  redirectTo: null,
  isLoaded: false
};

export default function(state=initialState, action) {
  switch(action.type) {
    case Actions.LOGIN_REQUIRED:
      return Object.assign({}, state, { redirectTo: action.payload });
    case Actions.LOGIN_SUCCESS:
    case Actions.CURRENT_USER_SUCCESS:
    case Actions.SIGNUP_SUCCESS:
      return Object.assign({}, state, action.payload, { isLoggedIn: true, isLoaded: true });
    case Actions.LOGOUT:
    case Actions.CURRENT_USER_FAILURE:
      return Object.assign({}, state, action.payload, { isLoaded: true });
  }
  return state;
}
