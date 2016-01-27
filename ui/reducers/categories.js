import _ from 'lodash';
import { Actions } from '../constants';

const initialState = {
  categoryMap: {},
  category: null,
};

export default function (state = initialState, action) {
  switch (action.type) {
    case Actions.CATEGORIES_LOADED:
      return Object.assign({}, state, {
        categoryMap: _.keyBy(action.payload, 'id'),
      });
    case Actions.GET_CATEGORY:
      return Object.assign({}, state, {
        category: state.categoryMap[action.payload],
      });
    case Actions.GET_CHANNEL_SUCCESS:
      return Object.assign({}, state, {
        categoryMap: Object.assign({}, state.categoryMap,
        _.keyBy(action.payload.categories, 'id')),
      });
    default:
      return state;
  }
}
