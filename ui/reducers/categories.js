import _ from 'lodash';
import { Actions } from '../constants';

const initialState = {
  categories: [],
  categoryMap: {},
  category: null,
};

export default function (state = initialState, action) {
  switch (action.type) {
    case Actions.CATEGORIES_LOADED:
      return Object.assign({}, state, {
        categories: action.payload,
        categoryMap: _.keyBy(action.payload, 'id'),
      });
    case Actions.GET_CATEGORY:
      return Object.assign({}, state, {
        category: state.categoryMap[action.payload],
      });
    default:
      return state;
  }
}
