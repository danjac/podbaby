import immutable from 'immutable';
import { Actions } from '../constants';

const initialState = immutable.Map({
  show: false,
});

export default function (state = initialState, action) {
  switch (action.type) {

    case Actions.OPEN_ADD_CHANNEL_FORM:
      return state.set('show', true);

    case Actions.ADD_CHANNEL_SUCCESS:
    case Actions.CLOSE_ADD_CHANNEL_FORM:
      return state.set('show', false);

    default:
      return state;
  }
}
