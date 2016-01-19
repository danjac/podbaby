import immutable from 'immutable';
import { Podcast } from '../records';
import { Actions } from '../constants';

const initialState = immutable.Map({
  podcasts: immutable.List(),
  showDetail: immutable.Set(),
  isLoading: false,
  page: immutable.Map({
    numPages: 0,
    numRows: 0,
    page: 1,
  }),
});

const podcastsFromJS = (payload) => {
  return immutable.List((payload || []).map(value => new Podcast(value)));
};


export default function (state = initialState, action) {
  switch (action.type) {

    case Actions.SHOW_PODCAST_DETAIL:

      return state.updateIn(['showDetail'], v => v.add(action.payload));

    case Actions.HIDE_PODCAST_DETAIL:

      return state.updateIn(['showDetail'], v => v.delete(action.payload));

    case Actions.CHANNEL_SEARCH_SUCCESS:
    case Actions.BOOKMARKS_SEARCH_SUCCESS:

      return state
        .set('isLoading', false)
        .set('podcasts', podcastsFromJS(action.payload));

    case Actions.GET_BOOKMARKS_SUCCESS:
    case Actions.GET_RECENT_PLAYS_SUCCESS:
    case Actions.LATEST_PODCASTS_SUCCESS:
    case Actions.GET_CHANNEL_SUCCESS:

      return state
        .set('isLoading', false)
        .set('page', immutable.fromJS(action.payload.page))
        .set('podcasts', podcastsFromJS(action.payload.podcasts));

    case Actions.SEARCH_SUCCESS:

      return state
        .set('isLoading', false)
        .set('podcasts', podcastsFromJS(action.payload.podcasts));

    case Actions.CLEAR_RECENT_PLAYS:

      return state.set('podcasts', immutable.List());

    case Actions.CLEAR_SEARCH:
    case Actions.SEARCH_REQUEST:
    case Actions.PODCASTS_REQUEST:

      return initialState.set('isLoading', true);

    case Actions.BOOKMARKS_SEARCH_FAILURE:
    case Actions.CHANNEL_SEARCH_FAILURE:
    case Actions.GET_BOOKMARKS_FAILURE:
    case Actions.GET_RECENT_PLAYS_FAILURE:
    case Actions.LATEST_PODCASTS_FAILURE:

      return initialState.set('isLoading', false);

    default:

      return state;
  }
}
