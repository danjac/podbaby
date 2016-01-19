import immutable from 'immutable';
import { Actions } from '../constants';

const initialState = immutable.Map({
  channel: null,
  query: '',
  isLoading: false,
});

export default function (state = initialState, action) {
  switch (action.type) {

    case Actions.CHANNEL_SEARCH_REQUEST:
      return state.set('query', action.payload);

    case Actions.ADD_CHANNEL_SUCCESS:
      return state
        .set('channel', action.payload)
        .set('isLoading', false);

    case Actions.GET_CHANNEL_SUCCESS:
      return state
        .set('query', '')
        .set('channel', action.payload)
        .set('isLoading', false);

    case Actions.GET_CHANNEL_FAILURE:
      return state
        .set('channel', null)
        .set('isLoading', false);

    case Actions.GET_CHANNEL_REQUEST:
      return state
        .set('channel', null)
        .set('isLoading', true);

    default:
      return state;
  }
}
