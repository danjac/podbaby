import { Actions } from '../constants';


const initialState = {
  podcasts: [],
  page: {
    numPages: 0,
    numRows: 0,
    page: 1
  }
};

export default function(state=initialState, action) {
  let podcasts;
  switch(action.type) {

    case Actions.SUBSCRIBE:
    case Actions.UNSUBSCRIBE:
      podcasts = state.podcasts.map(podcast => {
        if (podcast.channelId === action.payload) {
          podcast.isSubscribed = action.type === Actions.SUBSCRIBE;
        }
        return podcast;
      });
      return Object.assign({}, state, { podcasts });

    case Actions.GET_BOOKMARKS_SUCCESS:
      return action.payload;

    case Actions.GET_BOOKMARKS_FAILURE:
      return initialState;

    case Actions.DELETE_BOOKMARK:
      podcasts = _.reject(state.podcasts, bookmark => {
          return bookmark.id === action.payload;
      });
      return Object.assign({}, state, { podcasts });
  }
  return state;

}
