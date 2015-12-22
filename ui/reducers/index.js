import { combineReducers } from 'redux';
import { routeReducer } from 'redux-simple-router';

import authReducer from './auth';
import addChannelReducer from './add_channel';

export default combineReducers({
  routing: routeReducer,
  auth: authReducer,
  addChannel: addChannelReducer
});
