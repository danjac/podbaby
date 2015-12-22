import { Actions } from '../constants';

const initialState = {
  q: ""
};

export default function (state=initialState, action) {
  switch(action.type) {
    case Actions.SEARCH:
      return Object.assign({}, state, { q: action.payload });
  }
  return state;

}
