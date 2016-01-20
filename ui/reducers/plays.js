import { Actions } from '../constants';

const initialState = [];

export default function (state = initialState, action) {
  switch (action.type) {

    case Actions.CURRENT_USER:
    case Actions.LOGIN_SUCCESS:
      return action.payload && action.payload.plays ? action.payload.plays : initialState;

    case Actions.CLEAR_RECENT_PLAYS:
    case Actions.LOGOUT:
      return initialState;

    case Actions.CURRENTLY_PLAYING:

      return action.payload ?
        state.concat({
          podcastId: action.payload.id,
          createdAt: new Date(),
        }) : state;

    default:
      return state;
  }
}
