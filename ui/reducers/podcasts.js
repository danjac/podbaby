import _ from 'lodash';
import { Actions } from '../constants';

const initialState = {
  podcasts: [],
  showDetail: [],
  isLoading: true,
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
      showDetail = _.reject(state.showDetail, id => id === action.payload);
      return Object.assign({}, state, { showDetail });

    case Actions.ADD_BOOKMARK:
    case Actions.DELETE_BOOKMARK:

      podcasts = (state.podcasts || []).map(podcast => {
        if (podcast.id === action.payload) {
          podcast.isBookmarked = action.type === Actions.ADD_BOOKMARK;
        }
        return podcast;
      });
      return Object.assign({}, state, { podcasts });

    case Actions.SUBSCRIBE:
    case Actions.UNSUBSCRIBE:

      podcasts = (state.podcasts || []).map(podcast => {
        if (podcast.channelId === action.payload) {
          podcast.isSubscribed = action.type === Actions.SUBSCRIBE;
        }
        return podcast;
      });
      return Object.assign({}, state, { podcasts });

    case Actions.GET_BOOKMARKS_SUCCESS:
    case Actions.GET_RECENT_PLAYS_SUCCESS:
    case Actions.LATEST_PODCASTS_SUCCESS:
    case Actions.GET_CHANNEL_SUCCESS:

      return Object.assign({}, state, {
        podcasts: action.payload.podcasts,
        page: action.payload.page,
        isLoading: false
      });

    case Actions.SEARCH_SUCCESS:
      return Object.assign({}, state, {
        podcasts: action.payload.podcasts,
        isLoading: false
      });

    case Actions.CLEAR_RECENT_PLAYS:
      return Object.assign({}, state, { podcasts: [] });

    case Actions.CLEAR_SEARCH:

    case Actions.SEARCH_REQUEST:
    case Actions.PODCASTS_REQUEST:

    case Actions.GET_BOOKMARKS_FAILURE:
    case Actions.GET_RECENT_PLAYS_FAILURE:
    case Actions.LATEST_PODCASTS_FAILURE:

      return initialState;
  }
  return state;
}
