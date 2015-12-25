import { combineReducers } from 'redux';
import { routeReducer } from 'redux-simple-router';

import authReducer from './auth';
import searchReducer from './search';
import addChannelReducer from './add_channel';
import latestReducer from './latest';
import playerReducer from './player';

export default combineReducers({
  routing: routeReducer,
  auth: authReducer,
  search: searchReducer,
  addChannel: addChannelReducer,
  latest: latestReducer,
  player: playerReducer
});
