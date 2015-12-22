import React from 'react';
import { Router, Route, IndexRoute } from 'react-router';

import App from '../components/app';
import Front from '../components/front';
import Login from '../components/login';
import Signup from '../components/signup';
import EpisodeList from '../components/episode_list';
import PodcastList from '../components/podcast_list';
import PageNotFound from '../components/not_found';

import { loginRequired } from '../actions/auth';

export default function(store, history) {

    const requireAuth = (nextState, replaceState) => {
      const { auth } = store.getState();
      if (!auth.isLoggedIn) {
        store.dispatch(loginRequired(nextState.location.pathname));
        replaceState(null, "/login/");
      }
    };

    return (
      <Router history={history}>
        <Route path="/" component={App}>
          <IndexRoute component={Front} />
          <Route path="/podcasts/" onEnter={requireAuth}>
            <IndexRoute component={EpisodeList} />
            <Route path="/list/" component={PodcastList} />
          </Route>
          <Route path="/login/" component={Login} />
          <Route path="/signup/" component={Signup} />
          <Route path="*" component={PageNotFound} />
        </Route>
      </Router>
    );
}
