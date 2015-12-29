import _ from 'lodash';
import { Actions } from '../constants';

const initialState = {
  podcasts: [],
  showDetail: [],
  page: {
    numPages: 0,
    numRows: 0,
    page: 1
  }
};


export default function(state=initialState, action) {

  let podcasts, showDetail;

  switch(action.type) {

    case Actions.SHOW_PODCAST_DETAIL:
      showDetail = state.showDetail.concat(action.payload);
      return Object.assign({}, state, { showDetail });

    case Actions.HIDE_PODCAST_DETAIL:
      showDetail = _.reject(id => id === action.payload);
      return Object.assign({}, state, { showDetail });

    case Actions.ADD_BOOKMARK:
    case Actions.DELETE_BOOKMARK:
      podcasts = state.podcasts.map(podcast => {
        if (podcast.id === action.payload) {
          podcast.isBookmarked = action.type === Actions.ADD_BOOKMARK;
        }
        return podcast;
      });
      return Object.assign({}, state, { podcasts });

    case Actions.SUBSCRIBE:
    case Actions.UNSUBSCRIBE:
      podcasts = state.podcasts.map(podcast => {
        if (podcast.channelId === action.payload) {
          podcast.isSubscribed = action.type === Actions.SUBSCRIBE;
        }
        return podcast;
      });
      return Object.assign({}, state, { podcasts });

    case Actions.GET_CHANNEL_SUCCESS:
      return Object.assign({}, state, { podcasts: action.payload.podcasts, page: action.payload.page });

    case Actions.SEARCH_SUCCESS:
      return Object.assign({}, state, { podcasts: action.payload.podcasts });

    case Actions.GET_BOOKMARKS_SUCCESS:
    case Actions.LATEST_PODCASTS_SUCCESS:
      return Object.assign({}, state, { podcasts: action.payload.podcasts, page: action.payload.page });
      
    case Actions.GET_BOOKMARKS_FAILURE:
    case Actions.LATEST_PODCASTS_FAILURE:

      return initialState;
  }
  return state;
}
