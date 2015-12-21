import React from 'react';
import ReactDOM from 'react-dom';
import createHashHistory from 'history/lib/createHashHistory';
import { Provider } from 'react-redux';
import { syncReduxAndRouter } from 'redux-simple-router';

import Routes from './components/routes';
import configureStore from './store';

const history = createHashHistory();
const store = configureStore();

syncReduxAndRouter(history, store);

const Container = props => {
  return (
    <Provider store={store}>
      <Routes history={history} />
    </Provider>
  );
};

ReactDOM.render(<Container />, document.getElementById("app"));
