import { Actions } from '../constants';

const initialState = {
  channels: [],
  requestedChannels: [],
  filter: null,
  isLoading: false
};

export default function(state=initialState, action) {
let channels, filter;

  switch(action.type) {
    case Actions.UNSUBSCRIBE:

      channels = _.reject(
        state.channels,
        channel => channel.id === action.payload);

      return Object.assign({}, state, { channels, requestedChannels: channels });

    case Actions.FILTER_CHANNELS:
      filter = action.payload ? new RegExp(action.payload, "i") : null;
      if (filter) {
        channels = state.requestedChannels.filter(channel => channel.title.match(filter));
      } else {
        channels = state.requestedChannels.slice();
      }
      return Object.assign({}, state, { channels, filter });

    case Actions.GET_CHANNELS_REQUEST:
      return Object.assign({}, state, { isLoading: true });

    case Actions.GET_CHANNELS_SUCCESS:
      channels = action.payload || [];
      return Object.assign({}, state, {
        channels: channels,
        requestedChannels: channels,
        isLoading: false
      });

    case Actions.GET_CHANNELS_FAILURE:
      return Object.assign({}, state, {
        channels: [],
        requestedChannels: [],
        isLoading: false
      });

}

  return state;
}
