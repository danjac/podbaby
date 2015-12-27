import { combineReducers } from 'redux';
import { routeReducer } from 'redux-simple-router';

import authReducer from './auth';
import searchReducer from './search';
import addChannelReducer from './add_channel';
import latestReducer from './latest';
import playerReducer from './player';
import alertsReducer from './alerts';
import channelsReducer from './channels';
import channelReducer from './channel';

export default combineReducers({
  routing: routeReducer,
  auth: authReducer,
  search: searchReducer,
  addChannel: addChannelReducer,
  latest: latestReducer,
  player: playerReducer,
  alerts: alertsReducer,
  channels: channelsReducer,
  channel: channelReducer
});
