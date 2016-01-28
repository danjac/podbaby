import _ from 'lodash';
import { Actions } from '../constants';

const initialState = {
  podcasts: [],
  showDetail: [],
  isLoading: true,
  podcastCache: {},
  page: {
    numPages: 0,
    numRows: 0,
    page: 1,
  },
};


export default function (state = initialState, action) {
  let showDetail;

  switch (action.type) {

    case Actions.SHOW_PODCAST_DETAIL:
      showDetail = state.showDetail.concat(action.payload);
      return Object.assign({}, state, { showDetail });

    case Actions.HIDE_PODCAST_DETAIL:
      showDetail = _.reject(state.showDetail, id => id === action.payload);
      return Object.assign({}, state, { showDetail });

    case Actions.CHANNEL_SEARCH_SUCCESS:
    case Actions.BOOKMARKS_SEARCH_SUCCESS:
      return Object.assign({}, state, {
        podcasts: action.payload,
        podcastCache: Object.assign({}, state.podcastCache, _.keyBy(action.payload, 'id')),
        isLoading: false,
      });

    case Actions.GET_BOOKMARKS_SUCCESS:
    case Actions.GET_RECENT_PLAYS_SUCCESS:
    case Actions.LATEST_PODCASTS_SUCCESS:
    case Actions.GET_CHANNEL_SUCCESS:

      return Object.assign({}, state, {
        podcasts: action.payload.podcasts,
        page: action.payload.page,
        podcastCache: Object.assign({}, state.podcastCache, _.keyBy(action.payload.podcasts, 'id')),
        isLoading: false,
      });

    case Actions.SEARCH_SUCCESS:
      return Object.assign({}, state, {
        podcastCache: Object.assign({}, state.podcastCache, _.keyBy(action.payload.podcasts, 'id')),
        podcasts: action.payload.podcasts,
        isLoading: false,
      });

    case Actions.CLEAR_RECENT_PLAYS:
      return Object.assign({}, state, { podcasts: [] });

    case Actions.UPDATE_PODCAST_CACHE:
    case Actions.GET_PODCAST_SUCCESS:
      return Object.assign({}, state, {
        podcastCache: Object.assign({},
        state.podcastCache, { [action.payload.id]: action.payload }
        ),
      });

    case Actions.CLEAR_SEARCH:
    case Actions.SEARCH_REQUEST:
    case Actions.PODCASTS_REQUEST:
      return Object.assign({}, initialState, { podcastCache: state.podcastCache });

    case Actions.BOOKMARKS_SEARCH_FAILURE:
    case Actions.CHANNEL_SEARCH_FAILURE:

    case Actions.GET_BOOKMARKS_FAILURE:
    case Actions.GET_RECENT_PLAYS_FAILURE:
    case Actions.LATEST_PODCASTS_FAILURE:

      return Object.assign({}, initialState, { isLoading: false });
    default:

      return state;
  }
}
