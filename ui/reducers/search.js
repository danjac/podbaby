import { Actions } from '../constants';

const initialState = {
  query: "",
  channels: [],
  podcasts: []
};

export default function (state=initialState, action) {
  switch(action.type) {
    case Actions.SUBSCRIBE:
    case Actions.UNSUBSCRIBE:
      return Object.assign({}, state, {
        channels: state.channels.map(channel => {
          if (channel.id === action.payload) {
            channel.isSubscribed = action.type === Actions.SUBSCRIBE;
          }
          return channel;
        }),
        podcasts: state.podcasts.map(podcast => {
          if (podcast.channelId === action.payload) {
            podcast.isSubscribed = action.type === Actions.SUBSCRIBE;
          }
          return podcast;
        })
      });
    case Actions.ADD_BOOKMARK:
    case Actions.DELETE_BOOKMARK:
      return Object.assign({}, state, {
        podcasts: state.podcasts.map(podcast => {
          if (podcast.id === action.payload) {
            podcast.isBookmarked = action.type === Actions.ADD_BOOKMARK;
          }
          return podcast;
        })
      });
    case Actions.SEARCH:
      return Object.assign({}, state, { query: action.payload });
    case Actions.SEARCH_SUCCESS:
      let { channels, podcasts } = action.payload;
      return Object.assign({}, state, {
        channels: channels || [],
        podcasts: podcasts || []
      });
    case Actions.SEARCH_FAILURE:
      return initialState;
  }
  return state;

}
