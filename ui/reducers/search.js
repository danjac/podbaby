import { Actions } from '../constants';

const initialState = {
  query: "",
  channels: [],
  podcasts: [],
  numResults: 0
};

export default function (state=initialState, action) {
  switch(action.type) {
    case Actions.SUBSCRIBE:
    case Actions.UNSUBSCRIBE:
      return Object.assign({}, state, {
        channels: state.channels.map(channel => {
          if (channel.id === action.payload) {
            channel.isSubscribed = action.type === Actions.SUBSCRIBE;
          }
          return channel;

        })
      })
    case Actions.SEARCH:
      return Object.assign({}, state, { query: action.payload });
    case Actions.SEARCH_SUCCESS:
      let { channels, podcasts, numResults } = action.payload;
      return Object.assign({}, state, {
        channels: channels || [],
        podcasts: podcasts || [],
        numResults
      });
    case Actions.SEARCH_FAILURE:
      return initialState;
  }
  return state;

}
