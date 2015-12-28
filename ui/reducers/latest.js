import _ from 'lodash';
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

    case Actions.ADD_BOOKMARK:
    case Actions.DELETE_BOOKMARK:
      podcasts = state.podcasts.map(podcast => {
        if (podcast.id === action.payload) {
          podcast.isBookmarked = action.type === Actions.ADD_BOOKMARK;
        }
        return podcast;
      });
      return Object.assign({}, state, { podcasts });

    case Actions.UNSUBSCRIBE:
      podcasts = _.reject(state.podcasts, podcast => podcast.channelId === action.payload);
      return Object.assign({}, state, { podcasts });

    case Actions.LATEST_PODCASTS_SUCCESS:
      return action.payload;

    case Actions.LATEST_PODCASTS_FAILURE:
      return initialState;
  }
  return state;
}
