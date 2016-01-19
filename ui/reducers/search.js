import immutable from 'immutable';
import { Actions } from '../constants';

const initialState = immutable.Map({
  query: '',
});

export default function (state = initialState, action) {
  switch (action.type) {
    case Actions.CLEAR_SEARCH:
      return initialState;

    case Actions.SEARCH_REQUEST:
      return state.set('query', action.payload);

    default:
      return state;

  }
}
