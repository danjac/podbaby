import { Actions } from '../constants';

const initialState = {
  podcast: null,
  isPlaying: false,
  autoPlay: false,
  currentTime: 0,
};

export default function (state = initialState, action) {
  switch (action.type) {
    case Actions.CURRENTLY_PLAYING:
      return Object.assign({}, state, {
        podcast: action.payload,
        isPlaying: !!action.payload,
        currentTime: 0,
      });
    case Actions.ADD_BOOKMARK:
    case Actions.DELETE_BOOKMARK:
      if (state.podcast && state.podcast.id === action.payload) {
        const isBookmarked = action.type === Actions.ADD_BOOKMARK;
        const podcast = Object.assign({}, state.podcast, { isBookmarked });
        return Object.assign({}, state, { podcast });
      }
      return state;
    case Actions.PLAYER_TIME_UPDATE:
      return Object.assign({}, state, { currentTime: action.payload });
    case Actions.TOGGLE_AUTO_PLAY:
      return Object.assign({}, state, { autoPlay: !state.autoPlay });
    case Actions.CLOSE_PLAYER:
      return initialState;
    case Actions.RELOAD_PLAYER:
      return action.payload || initialState;
    default:
      return state;
  }
}
