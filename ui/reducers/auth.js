import { Actions } from '../constants';

const initialState = {
  isLoggedIn: false,
  name: null,
  redirectTo: null
};

export default function(state=initialState, action) {
  switch(action.type) {
    case Actions.LOGIN_REQUIRED:
      return Object.assign({}, state, { redirectTo: action.payload });
    case Actions.LOGIN_SUCCESS:
    case Actions.SIGNUP_SUCCESS:
      return Object.assign({}, state, action.payload, { isLoggedIn: true });
    case Actions.LOGOUT:
      return initialState;
  }
  return state;
}
