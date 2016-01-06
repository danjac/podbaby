import { Actions } from '../constants';

const initialState = {
  query: ""
};

export default function(state=initialState, action) {
  switch(action.type) {
    case Actions.BOOKMARKS_SEARCH_REQUEST:
      return Object.assign({}, state, { query: action.payload });
    case Actions.CLEAR_BOOKMARKS_SEARCH:
      return Object.assign({}, state, { query: "" });
  }
  return state;
};
