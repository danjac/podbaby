import React from 'react';
import ReactDOM from 'react-dom';
import createHashHistory from 'history/lib/createHashHistory';
import { syncReduxAndRouter } from 'redux-simple-router';

import Root from './containers/root';
import { auth, player, categories } from './actions';

import configureStore from './store';
import configureRoutes from './routes';

const history = createHashHistory();
const store = configureStore();
const routes = configureRoutes(store, history);

syncReduxAndRouter(history, store);

store.dispatch(auth.setCurrentUser(window.__DATA__.user));
store.dispatch(categories.loadCategories(window.__DATA__.categories));
store.dispatch(player.reloadPlayer());

ReactDOM.render(
  <Root store={store} routes={routes} />,
  document.getElementById('app'));
