import immutable from 'immutable';
import { Podcast } from '../records';
import { Actions } from '../constants';

const initialState = immutable.Map({
  podcast: null,
  isPlaying: false,
  currentTime: 0,
});

export default function (state = initialState, action) {
  switch (action.type) {

    case Actions.CURRENTLY_PLAYING:
      return state
        .set('podcast', new Podcast(action.payload))
        .set('isPlaying', true)
        .set('currentTime', 0);

    case Actions.ADD_BOOKMARK:
    case Actions.DELETE_BOOKMARK:

      return state.updateIn(['podcast'], podcast => {
        if (podcast && podcast.get('id') === action.payload) {
          return podcast.set('isBookmarked', action.type === Actions.ADD_BOOKMARK);
        }
        return podcast;
      });

    case Actions.PLAYER_TIME_UPDATE:
      return state.set('currentTime', action.payload);

    case Actions.CLOSE_PLAYER:
      return initialState;

    case Actions.RELOAD_PLAYER:
      // reload properly
      return action.payload ? immutable.Map(immutable.fromJS(action.payload, (k, v) => {
        return k === 'podcast' ? new Podcast(v) : v;
      })) : initialState;

    default:
      return state;
  }
}
