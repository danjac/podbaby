import { Actions } from '../constants';

const initialState = {
  channel: null,
  isLoading: false
};

export default function(state=initialState, action) {

  switch(action.type) {

    case Actions.SUBSCRIBE:
    case Actions.UNSUBSCRIBE:
      if (state.channel && state.channel.id === action.payload) {
        let channel = Object.assign({}, state.channel, { isSubscribed: action.type === Actions.SUBSCRIBE });
        return Object.assign({}, state, { channel });
      }
      return state;

    case Actions.GET_CHANNEL_SUCCESS:
      return Object.assign({}, state, { channel: action.payload.channel, isLoading: false });

    case Actions.GET_CHANNEL_FAILURE:
      return Object.assign({}, state, { channel: null, isLoading: false });

    case Actions.GET_CHANNEL_REQUEST:
      return Object.assign({}, state, { channel: null, isLoading: true });
  }
  return state;
}
