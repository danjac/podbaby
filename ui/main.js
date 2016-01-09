import React from 'react';
import ReactDOM from 'react-dom';
import createHashHistory from 'history/lib/createHashHistory';
import { Provider } from 'react-redux';
import { syncReduxAndRouter } from 'redux-simple-router';
import DocumentTitle from 'react-document-title';

import { auth, player } from './actions';
import routes from './routes';
import configureStore from './store';
import { getTitle } from './containers/utils';

const history = createHashHistory();
const store = configureStore();

syncReduxAndRouter(history, store);

const Container = props => {

  return (
  <DocumentTitle title={getTitle()}>
    <Provider store={store}>
      {routes(store, history)}
    </Provider>
  </DocumentTitle>
  );
};

store.dispatch(auth.setCurrentUser(window.user));
store.dispatch(player.reloadPlayer());

ReactDOM.render(<Container />, document.getElementById("app"));
