import { combineReducers } from 'redux';
import { routeReducer } from 'redux-simple-router';

import authReducer from './auth';

export default combineReducers({
  routing: routeReducer,
  auth: authReducer
});
