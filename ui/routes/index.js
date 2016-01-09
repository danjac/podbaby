import React from 'react';

import {
  Router,
  Route,
  IndexRoute,
  IndexRedirect
} from 'react-router';

import App from '../containers/app';
import Front from '../containers/front';
import Login from '../containers/login';
import Signup from '../containers/signup';
import Search from '../containers/search';
import Latest from '../containers/latest';
import Recent from '../containers/recent';
import Subscriptions from '../containers/subscriptions';
import Bookmarks from '../containers/bookmarks';
import Channel from '../containers/channel';
import Podcast from '../containers/podcast';
import User from '../containers/user';
import PageNotFound from '../containers/not_found';

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
      replaceState(null, "/new/");
    } else {
      replaceState(null, "/front/");
    }
  };

  return (
    <Router history={history}>
      <Route path="/" component={App}>
        <IndexRoute onEnter={redirectIfLoggedIn} />

        <Route component={Front} path="/front/" />
        <Route path="/member/" onEnter={requireAuth}>

          <IndexRedirect to="/member/subscriptions/" />

          <Route path="/member/subscriptions/"
                 component={Subscriptions}
                 onEnter={() => actions.channels.getChannels()} />

          <Route path="/member/bookmarks/"
                 component={Bookmarks}
                 onEnter={() => actions.bookmarks.getBookmarks()} />

          <Route path="/member/recent/"
                 component={Recent}
                 onEnter={() => actions.plays.getRecentlyPlayed()} />
        </Route>

        <Route path="/new/"
               component={Latest}
               onEnter={() => actions.latest.getLatestPodcasts()} />

        <Route path="/search/" component={Search} />

        <Route path="/channel/:id/"
               component={Channel}
               onEnter={nextState => actions.channel.getChannel(nextState.params.id)} />

        <Route path="/podcast/:id/"
               component={Podcast}
               onEnter={nextState => actions.podcasts.getPodcast(nextState.params.id)} />

        <Route path="/login/" component={Login} />
        <Route path="/signup/" component={Signup} />
        <Route path="/user/" component={User} onEnter={requireAuth} />
        <Route path="*" component={PageNotFound} />
      </Route>
    </Router>
  );
}
