import { Actions } from '../constants';

const initialState = {
  channel: null
};

export default function(state=initialState, action) {

  switch(action.type) {

    case Actions.SUBSCRIBE:
    case Actions.UNSUBSCRIBE:
      if (state.channel && state.channel.id === action.payload) {
        channel = Object.assign({}, state.channel, { isSubscribed: action.type === Actions.SUBSCRIBE });
        return Object.assign({}, state, { channel });
      }
      return state;

    case Actions.GET_CHANNEL_SUCCESS:
      return Object.assign({}, state, { channel: action.payload.channel });

    case Actions.GET_CHANNEL_FAILURE:
      return initialState;
  }
  return state;
}
