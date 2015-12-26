import { Actions } from '../constants';

const initialState = {
  podcast: null,
  isPlaying: false
};

export default function(state=initialState, action) {
  switch(action.type) {
    case Actions.CURRENTLY_PLAYING:
      return { podcast: action.payload, isPlaying: !!action.payload };
  }
  return state;
}
