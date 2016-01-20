import immutable from 'immutable';
import { Alert } from '../records';
import { Actions } from '../constants';

const initialState = immutable.List();

export default function (state = initialState, action) {
  switch (action.type) {
    case Actions.ADD_ALERT:
      return state.push(new Alert(action.payload));
    case Actions.DISMISS_ALERT:
      return state.filterNot(alert => alert.get('id') === action.payload);
    default:
      return state;
  }
}
