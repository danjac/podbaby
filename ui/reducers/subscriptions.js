import immutable from 'immutable';
import { Actions } from '../constants';

const initialState = immutable.Set();

export default function (state = initialState, action) {
  switch (action.type) {

    case Actions.LOGIN_SUCCESS:
    case Actions.CURRENT_USER:

      return immutable.Set(
        action.payload ? action.payload.subscriptions : []
      );

    case Actions.LOGOUT:
      return initialState;

    case Actions.ADD_CHANNEL_SUCCESS:
      return state.add(action.payload.id);

    case Actions.SUBSCRIBE:
      return state.add(action.payload);

    case Actions.UNSUBSCRIBE:
      return state.delete(action.payload);

    default:
      return state;
  }
}
