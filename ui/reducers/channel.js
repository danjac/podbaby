import { Actions } from '../constants';

const initialState = null;

export default function(state=initialState, action) {
  switch(action.type) {

    case Actions.SUBSCRIBE:
    case Actions.UNSUBSCRIBE:
      return state && state.id === action.payload ?
        Object.assign({}, state, { isSubscribed: action.type === Actions.SUBSCRIBE }) : state;

    case Actions.ADD_BOOKMARK:
    case Actions.DELETE_BOOKMARK:
      if (state && state.podcasts) {
        let podcasts = state.podcasts.map(podcast => {
          if (podcast.id === action.payload) {
            podcast.isBookmarked = action.type === Actions.ADD_BOOKMARK;
          }
          return podcast;
        });
        return Object.assign({}, state, { podcasts });
      }
      return state;

    case Actions.GET_CHANNEL_SUCCESS:
      return action.payload;

    case Actions.GET_CHANNEL_FAILURE:
      return initialState;
  }
  return state;
}
