import React from 'react';

import {
  Router,
  Route,
  IndexRoute,
  IndexRedirect
} from 'react-router';

import App from '../components/app';
import Front from '../components/front';
import Login from '../components/login';
import Signup from '../components/signup';
import Search from '../components/search';
import Latest from '../components/latest';
import Recent from '../components/recent';
import Subscriptions from '../components/subscriptions';
import Bookmarks from '../components/bookmarks';
import Channel from '../components/channel';
import User from '../components/user';
import PageNotFound from '../components/not_found';

import * as actionCreators from '../actions';
import { bindAllActionCreators } from '../actions/utils';

export default function(store, history) {

  const actions = bindAllActionCreators(actionCreators, store.dispatch);

  const requireAuth = (nextState, replaceState) => {
    const { auth } = store.getState();
    if (!auth.isLoggedIn) {
      actions.alerts.warning("You have to be signed in first");
      actions.auth.loginRequired(nextState.location.pathname);
      replaceState(null, "/login/");
    }
  };

  const redirectIfLoggedIn = (nextState, replaceState) => {
    const { auth } = store.getState();
    if (auth.isLoggedIn) {
      replaceState(null, "/podcasts/");
    }
  };

  return (
    <Router history={history}>
      <Route path="/" component={App}>
        <IndexRoute component={Front} onEnter={redirectIfLoggedIn} />
        <Route path="/podcasts/" onEnter={requireAuth}>
          <IndexRedirect to="/podcasts/new/" />

          <Route path="/podcasts/new/"
                 component={Latest}
                 onEnter={() => actions.latest.getLatestPodcasts()} />

          <Route path="/podcasts/search/" component={Search} />

          <Route path="/podcasts/subscriptions/"
                 component={Subscriptions}
                 onEnter={() => actions.channels.getChannels()} />

          <Route path="/podcasts/bookmarks/"
                 component={Bookmarks}
                 onEnter={() => actions.bookmarks.getBookmarks()} />

          <Route path="/podcasts/recent/"
                 component={Recent}
                 onEnter={() => actions.plays.getRecentlyPlayed()} />

          <Route path="/podcasts/channel/:id/"
                 component={Channel}
                 onEnter={nextState => actions.channel.getChannel(nextState.params.id)} />

        </Route>
        <Route path="/login/" component={Login} />
        <Route path="/signup/" component={Signup} />
        <Route path="/user/" component={User} onEnter={requireAuth} />
        <Route path="*" component={PageNotFound} />
      </Route>
    </Router>
  );
}
