import immutable from 'immutable';
import { Actions } from '../constants';

const initialState = immutable.Map({
  query: '',
  bookmarks: immutable.Set(),
});

export default function (state = initialState, action) {
  switch (action.type) {

    case Actions.BOOKMARKS_SEARCH_REQUEST:
      return state.set('query', action.payload);

    case Actions.LOGIN_SUCCESS:
    case Actions.CURRENT_USER:

      return state.set('bookmarks', immutable.Set(
        action.payload && action.payload.bookmarks ? action.payload.bookmarks : []
      ));

    case Actions.LOGOUT:
      return initialState;

    case Actions.ADD_BOOKMARK:
      return state.updateIn(['bookmarks'], v => v.add(action.payload));

    case Actions.DELETE_BOOKMARK:
      return state.updateIn(['bookmarks'], v => v.delete(action.payload));

    case Actions.CLEAR_BOOKMARKS_SEARCH:
      return state.set('query', '');

    default:
      return state;
  }
}
