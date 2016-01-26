import { applyMiddleware, createStore, compose } from 'redux';
import thunkMiddleware from 'redux-thunk';
import createLogger from 'redux-logger';

import { reduxRouterMiddleware, apiErrorMiddleware } from '../middleware';
import reducer from '../reducers';

const loggerMiddleware = createLogger();

const middleware = applyMiddleware(
  reduxRouterMiddleware,
  thunkMiddleware,
  apiErrorMiddleware,
  loggerMiddleware
);

const createFinalStore = compose(middleware)(createStore);

export default function configureStore(initialState) {
  return createFinalStore(reducer, initialState);
}
