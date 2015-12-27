
import { Actions } from '../constants';

const initialState = {
  show: false
};

export default function(state=initialState, action) {
  switch(action.type) {

    case Actions.OPEN_ADD_CHANNEL_FORM:
      return Object.assign({}, state, { show: true });

    case Actions.ADD_CHANNEL_SUCCESS:
    case Actions.ADD_CHANNEL_FAILURE:
    case Actions.CLOSE_ADD_CHANNEL_FORM:
      return Object.assign({}, state, { show: false });

  }
  return state;
}
