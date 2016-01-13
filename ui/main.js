import React from 'react';
import ReactDOM from 'react-dom';
import createHashHistory from 'history/lib/createHashHistory';
import { syncReduxAndRouter } from 'redux-simple-router';

import Root from './containers/root';
import { auth, player } from './actions';

import configureStore from './store';
import configureRoutes from './routes';

const history = createHashHistory();
const store = configureStore();
const routes = configureRoutes(store, history);

syncReduxAndRouter(history, store);

store.dispatch(auth.setCurrentUser(window.user));
store.dispatch(player.reloadPlayer());

ReactDOM.render(
  <Root store={store} routes={routes} />,
  document.getElementById("app"));
