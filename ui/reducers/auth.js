import immutable from 'immutable';
import { Actions } from '../constants';

const Auth = immutable.Record({
  name: null,
  email: null,
  redirectTo: null,
  showRecoverPasswordForm: false,
  isLoaded: true,
  isLoggedIn: false,
});

const initialState = new Auth();

export default function (state = initialState, action) {
  switch (action.type) {

    case Actions.LOGIN_REQUIRED:
      return state.set('redirectTo', action.payload);

    case Actions.LOGIN_SUCCESS:
    case Actions.SIGNUP_SUCCESS:
      return state
        .set('name', action.payload.name)
        .set('email', action.payload.email)
        .set('isLoggedIn', true);

    case Actions.CURRENT_USER:

      const { name, email } = action.payload ? action.payload : [null, null];
      const isLoggedIn = Boolean(name && email);
      return state
        .set('name', name)
        .set('email', email)
        .set('isLoggedIn', isLoggedIn);

    case Actions.CHANGE_EMAIL_SUCCESS:
      return state.set('email', action.payload);

    case Actions.OPEN_RECOVER_PASSWORD_FORM:
    case Actions.CLOSE_RECOVER_PASSWORD_FORM:

      return state
        .set('showRecoverPasswordForm', action.type === Actions.OPEN_RECOVER_PASSWORD_FORM);

    case Actions.DELETE_ACCOUNT_SUCCESS:
    case Actions.LOGOUT:
      return initialState;
    default:
      return state;
  }
}
