import { Actions } from '../constants';

const initialState = [];

export default function (state = initialState, action) {
  switch (action.type) {
    case Actions.CURRENT_USER:
    case Actions.LOGIN_SUCCESS:
      return action.payload && action.payload.plays ? action.payload.plays : initialState;
    case Actions.LOGOUT:
      return initialState;
    default:
      return state;
  }
}
