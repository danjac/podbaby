import { createStore, applyMiddleware, compose } from 'redux';
import thunkMiddleware from 'redux-thunk';
import createLogger from 'redux-logger';

import reducer from '../reducers';

const loggingMiddleware = createLogger({
  level: 'info',
  collapsed: true,
  logger: console
});

const createStoreWithMiddleware = compose(
  applyMiddleware(
    thunkMiddleware,
    loggingMiddleware
  ),
)(createStore);

export default function configureStore(initialState) {
  return createStoreWithMiddleware(reducer, initialState);
}
