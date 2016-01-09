import { combineReducers } from 'redux';
import { reducer as formReducer } from 'redux-form';
import { routeReducer } from 'redux-simple-router';

import authReducer from './auth';
import searchReducer from './search';
import addChannelReducer from './add_channel';
import playerReducer from './player';
import alertsReducer from './alerts';
import channelsReducer from './channels';
import channelReducer from './channel';
import podcastsReducer from './podcasts';
import podcastReducer from './podcast';
import bookmarksReducer from './bookmarks';
import subscriptionsReducer from './subscriptions';

export default combineReducers({
  routing: routeReducer,
  form: formReducer,
  auth: authReducer,
  search: searchReducer,
  addChannel: addChannelReducer,
  player: playerReducer,
  alerts: alertsReducer,
  channels: channelsReducer,
  channel: channelReducer,
  bookmarks: bookmarksReducer,
  subscriptions: subscriptionsReducer,
  podcasts: podcastsReducer,
  podcast: podcastReducer
});
