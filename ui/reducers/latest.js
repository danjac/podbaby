import _ from 'lodash';
import { Actions } from '../constants';

const initialState = [];

export default function(state=initialState, action) {
  switch(action.type) {
    case Actions.ADD_BOOKMARK:
    case Actions.DELETE_BOOKMARK:
      return state.map(podcast => {
        if (podcast.id === action.payload) {
          podcast.isBookmarked = action.type === Actions.ADD_BOOKMARK;
        }
        return podcast;
      });
    case Actions.UNSUBSCRIBE:
      return _.reject(state, podcast => podcast.channelId === action.payload);
    case Actions.LATEST_PODCASTS_SUCCESS:
      return action.payload || [];
    case Actions.LATEST_PODCASTS_FAILURE:
      return [];
  }
  return state;
}
