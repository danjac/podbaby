import _ from 'lodash';
import { Actions } from '../constants';

const initialState = [];

export default function (state = initialState, action) {
  switch (action.type) {
    case Actions.ADD_ALERT:
      return state.concat(action.payload);
    case Actions.DISMISS_ALERT:
      return _.reject(state, alert => alert.id === action.payload);
    default:
      return state;
  }
}
