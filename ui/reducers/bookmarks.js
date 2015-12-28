import { Actions } from '../constants';

const initialState = [];

export default function(state=initialState, action) {
  switch(action.type) {
    
    case Actions.GET_BOOKMARKS_SUCCESS:
      return action.payload;

    case Actions.DELETE_BOOKMARK:
      return _.reject(state, bookmark => {
          return bookmark.id === action.payload;
      });


  }
  return state;

}
