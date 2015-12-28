import React from 'react';
import { Router, Route, IndexRoute, IndexRedirect } from 'react-router';

import App from '../components/app';
import Front from '../components/front';
import Login from '../components/login';
import Signup from '../components/signup';
import Search from '../components/search';
import Latest from '../components/latest';
import Subscriptions from '../components/subscriptions';
import Bookmarks from '../components/bookmarks';
import Channel from '../components/channel';
import PageNotFound from '../components/not_found';

import { alerts } from '../actions';
import { loginRequired } from '../actions/auth';

export default function(store, history) {

    const requireAuth = (nextState, replaceState) => {
      const { auth } = store.getState();
      if (!auth.isLoggedIn) {
        store.dispatch(alerts.warning("You have to be signed in first"));
        store.dispatch(loginRequired(nextState.location.pathname));
        replaceState(null, "/login/");
      }
    };

    return (
      <Router history={history}>
        <Route path="/" component={App}>
          <IndexRoute component={Front} />
          <Route path="/podcasts/" onEnter={requireAuth}>
            <IndexRedirect to="/podcasts/new/" />
            <Route path="/podcasts/new/" component={Latest} />
            <Route path="/podcasts/search/" component={Search} />
            <Route path="/podcasts/subscriptions/" component={Subscriptions} />
            <Route path="/podcasts/bookmarks/" component={Bookmarks} />
            <Route path="/podcasts/channel/:id/" component={Channel} />
          </Route>
          <Route path="/login/" component={Login} />
          <Route path="/signup/" component={Signup} />
          <Route path="*" component={PageNotFound} />
        </Route>
      </Router>
    );
}
