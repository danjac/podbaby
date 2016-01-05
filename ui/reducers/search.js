import { Actions } from '../constants';

const initialState = {
  query: "",
  channels: []
};

export default function (state=initialState, action) {
  switch(action.type) {
    case Actions.SUBSCRIBE:
    case Actions.UNSUBSCRIBE:
      return Object.assign({}, state, {
        channels: (state.channels || []).map(channel => {
          if (channel.id === action.payload) {
            channel.isSubscribed = action.type === Actions.SUBSCRIBE;
          }
          return channel;
        }),
      });

    case Actions.CLEAR_SEARCH:
      return initialState;

    case Actions.SEARCH_REQUEST:
      return Object.assign({}, state, { query: action.payload });

    case Actions.SEARCH_SUCCESS:
      let { channels } = action.payload;
      return Object.assign({}, state, {
        channels,
      });
    case Actions.SEARCH_FAILURE:
      return initialState;
  }
  return state;

}
