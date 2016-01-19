import immutable from 'immutable';
import { Podcast } from '../records';
import { Actions } from '../constants';


const initialState = immutable.Map({
  podcast: null,
  isLoading: false,
});

export default function (state = initialState, action) {
  switch (action.type) {

    case Actions.GET_PODCAST_REQUEST:
      return state.set('isLoading', true);

    case Actions.GET_PODCAST_SUCCESS:
      return state
        .set('isLoading', false)
        .set('podcast', new Podcast(action.payload));

    case Actions.GET_PODCAST_FAILURE:
      return state
        .set('isLoading', false)
        .set('podcast', null);

    default:
      return state;
  }
}
