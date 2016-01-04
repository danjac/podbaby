import { Actions } from '../constants';

const initialState = {
  podcast: null,
  isPlaying: false,
  currentTime: 0
};

export default function(state=initialState, action) {
  switch(action.type) {
    case Actions.CURRENTLY_PLAYING:
      return Object.assign({}, state, {
        podcast: action.payload,
        isPlaying: !!action.payload
      });
    case Actions.PLAYER_TIME_UPDATE:
      return Object.assign({}, state, { currentTime: action.payload });
    case Actions.CLOSE_PLAYER:
      return initialState;
    case Actions.RELOAD_PLAYER:
      return action.payload || initialState;
  }
  return state;
}
