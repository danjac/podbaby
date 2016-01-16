import { Actions } from '../constants';

const initialState = {
  podcast: null,
  isLoading: false,
};

export default function (state = initialState, action) {
  switch (action.type) {
    case Actions.GET_PODCAST_REQUEST:
      return Object.assign({}, state, { isLoading: true });
    case Actions.GET_PODCAST_SUCCESS:
      return Object.assign({}, state, { isLoading: false, podcast: action.payload });
    case Actions.GET_PODCAST_FAILURE:
      return Object.assign({}, state, { isLoading: false, podcast: null });
    default:
      return state;
  }
}
