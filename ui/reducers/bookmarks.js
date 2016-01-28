import { Actions } from '../constants';

const initialState = {
  query: '',
  playing: null,
  bookmarks: [],
};

export default function (state = initialState, action) {
  let bookmarks;

  switch (action.type) {

    case Actions.BOOKMARKS_SEARCH_REQUEST:
      return Object.assign({}, state, { query: action.payload });

    case Actions.LOGIN_SUCCESS:
    case Actions.CURRENT_USER:
      bookmarks = (action.payload && action.payload.bookmarks) || [];
      return Object.assign({}, state, { bookmarks });

    case Actions.LOGOUT:
      return Object.assign({}, state, { bookmarks: [] });

    case Actions.ADD_BOOKMARK:
      bookmarks = state.bookmarks.slice();
      bookmarks.unshift(action.payload);
      return Object.assign({}, state, { bookmarks });

    case Actions.DELETE_BOOKMARK:
      let playing = state.playing;
      if (state.bookmarks && state.playing === action.payload) {
        const prev = state.bookmarks.indexOf(state.playing) - 1;
        playing = state.bookmarks[prev < 1 ? 0 : prev];
      }
      bookmarks = state.bookmarks.filter(id => id !== action.payload);
      return Object.assign({}, state, { bookmarks, playing });

    case Actions.CLEAR_BOOKMARKS_SEARCH:
      return Object.assign({}, state, { query: '' });

    case Actions.BOOKMARKS_CURRENTLY_PLAYING:
      return Object.assign({}, state, { playing: action.payload });

    default:
      return state;
  }
}
