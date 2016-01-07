import { Actions } from '../constants';

const initialState = [];

export default function(state=initialState, action) {

  switch(action.type) {

    case Actions.LOGIN_SUCCESS:
    case Actions.CURRENT_USER:
      return (action.payload && action.payload.subscriptions) || [];

    case Actions.LOGOUT:
      return [];

    case Actions.SUBSCRIBE:
      return state.concat(action.payload);

    case Actions.UNSUBSCRIBE:
      return state.filter(id => id !== action.payload);
  }

  return state;

}

