import { Actions } from '../constants';

const initialState = {
  channel: null,
  podcasts: [],
  page: {
    numPages: 0,
    numRows: 0,
    page: 1
  }
};

export default function(state=initialState, action) {

  let channel, podcasts;

  switch(action.type) {

    case Actions.SUBSCRIBE:
    case Actions.UNSUBSCRIBE:
      if (state.channel && state.channel.id === action.payload) {
        channel = Object.assign({}, state.channel, { isSubscribed: action.type === Actions.SUBSCRIBE });
        return Object.assign({}, state, { channel });
      }
      return state;

    case Actions.ADD_BOOKMARK:
    case Actions.DELETE_BOOKMARK:
      podcasts = state.podcasts.map(podcast => {
        if (podcast.id === action.payload) {
          podcast.isBookmarked = action.type === Actions.ADD_BOOKMARK;
        }
      });
      return Object.assign({}, state, { podcasts });

    case Actions.GET_CHANNEL_SUCCESS:
      return action.payload;

    case Actions.GET_CHANNEL_FAILURE:
      return initialState;
  }
  return state;
}
