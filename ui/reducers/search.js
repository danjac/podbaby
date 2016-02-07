import { Actions } from '../constants';

const initialState = {
  query: '',
  type: 'podcasts',
};

export default function (state = initialState, action) {
  switch (action.type) {
    case Actions.CLEAR_SEARCH:
      return initialState;

    case Actions.SEARCH_REQUEST:
      return action.payload;

    default:
      return state;

  }
}
