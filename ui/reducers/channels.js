import immutable from 'immutable';
import { Channel } from '../records';
import { Actions } from '../constants';

const initialState = immutable.Map({
  channels: immutable.List(),
  filter: '',
  isLoading: false,
  page: 1,
});

const channelsFromJS = payload => {
  return immutable.List((payload || []).map(v => new Channel(v)));
};

export default function (state = initialState, action) {
  switch (action.type) {

    case Actions.FILTER_CHANNELS:
      return state
        .set('filter', action.payload)
        .set('page', 1);

    case Actions.SELECT_CHANNELS_PAGE:
      return state.set('page', action.payload);

    case Actions.GET_CHANNELS_REQUEST:
      return state.set('isLoading', true);

    case Actions.SEARCH_SUCCESS:

      return state
        .set('channels', channelsFromJS(action.payload.channels))
        .set('isLoading', false)
        .set('filter', '');

    case Actions.GET_CHANNELS_SUCCESS:

      return state
        .set('channels', channelsFromJS(action.payload))
        .set('page', 1)
        .set('isLoading', false)
        .set('filter', '');

    case Actions.CLEAR_SEARCH:
    case Actions.SEARCH_FAILURE:
    case Actions.GET_CHANNELS_FAILURE:
      return initialState;

    default:
      return state;
  }
}
