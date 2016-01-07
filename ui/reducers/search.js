import { Actions } from '../constants';

const initialState = {
  query: ""
};

export default function (state=initialState, action) {
  switch(action.type) {

    case Actions.CLEAR_SEARCH:
      return initialState;

    case Actions.SEARCH_REQUEST:
      return Object.assign({}, state, { query: action.payload });

  }
  return state;

}
