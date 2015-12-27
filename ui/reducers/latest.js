import { Actions } from '../constants';

const initialState = [];

export default function(state=initialState, action) {
  switch(action.type) {
    case Actions.LATEST_PODCASTS_SUCCESS:
      return action.payload || [];
    case Actions.LATEST_PODCASTS_FAILURE:
      return [];
  }
  return state;
}
