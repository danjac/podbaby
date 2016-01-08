import { Actions } from '../constants';

const initialState = {
  channels: [],
  filter: "",
  isLoading: false,
  page: 1
};

export default function(state=initialState, action) {

  let channels;

  switch(action.type) {

    case Actions.FILTER_CHANNELS:
      return Object.assign({}, state, { filter: action.payload, page: 1 });

    case Actions.SELECT_CHANNELS_PAGE:
      return Object.assign({}, state, { page: action.payload });

    case Actions.GET_CHANNELS_REQUEST:
      return Object.assign({}, state, { isLoading: true });

    case Actions.SEARCH_SUCCESS:

      channels = action.payload.channels || [];
      return Object.assign({}, state, {
        channels: channels,
        isLoading: false,
        filter: ""
      });

    case Actions.GET_CHANNELS_SUCCESS:
      channels = action.payload || [];
      return Object.assign({}, state, {
        channels: channels,
        isLoading: false,
        filter: "",
        page: 1
      });

    case Actions.CLEAR_SEARCH:
    case Actions.SEARCH_FAILURE:
    case Actions.GET_CHANNELS_FAILURE:
      return Object.assign({}, state, {
        channels: [],
        isLoading: false,
        filter: "",
        page: 1
      });

}

  return state;
}
