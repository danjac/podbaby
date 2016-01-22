import { Actions } from '../constants';

const initialState = [];

export default function (state = initialState, action) {
  switch (action.type) {
    case Actions.CATEGORIES_LOADED:
      return action.payload;
    default:
      return state;
  }
}
