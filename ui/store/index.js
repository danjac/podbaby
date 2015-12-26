import { createStore, applyMiddleware, compose } from 'redux';
import thunkMiddleware from 'redux-thunk';
import createLogger from 'redux-logger';

import { pushPath } from 'redux-simple-router';

import { alerts } from '../actions';
import reducer from '../reducers';

const loggingMiddleware = createLogger({
  level: 'info',
  collapsed: true,
  logger: console
});

// should catch any API errors and act accordingly
const apiErrorMiddleware = store => next => action => {
    let result = next(action);

    if (result.payload && result.payload.error) {

      const { error } = result.payload;

      switch(error.status) {
        case 401:
          store.dispatch(alerts.warning('You must be logged in to continue'))
          store.dispatch(pushPath("/login/"));
          break;
      }

    }
    return result;
};

const createStoreWithMiddleware = compose(
  applyMiddleware(
    thunkMiddleware,
    apiErrorMiddleware,
    loggingMiddleware,
  ),
)(createStore);

export default function configureStore(initialState) {
  return createStoreWithMiddleware(reducer, initialState);
}
