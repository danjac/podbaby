import { Actions } from '../constants';

const initialState = {
  channels: [],
  requestedChannels: [],
  filter: null,
  isLoading: false
};

const subscribeChannels = (channels, channel_id, subscribed) => {
  return channels.map(channel => {
    if (channel.id === channel_id) {
      channel.isSubscribed = subscribed;
    }
    return channel;
  });
}

export default function(state=initialState, action) {

  let channels, requestedChannels, filter, subscribed;

  switch(action.type) {

    case Actions.SUBSCRIBE:
    case Actions.UNSUBSCRIBE:

      subscribed = action.type === Actions.SUBSCRIBE;
      channels = subscribeChannels(state.channels, action.payload, subscribed);
      requestedChannels = subscribeChannels(state.requestedChannels, action.payload, subscribed);

      return Object.assign({}, state, { channels, requestedChannels });

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
