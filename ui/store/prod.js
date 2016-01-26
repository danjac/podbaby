import { applyMiddleware, createStore, compose } from 'redux';
import thunkMiddleware from 'redux-thunk';

import { reduxRouterMiddleware, apiErrorMiddleware } from '../middleware';
import reducer from '../reducers';

const middleware = applyMiddleware(
  reduxRouterMiddleware,
  thunkMiddleware,
  apiErrorMiddleware
);

const createFinalStore = compose(middleware)(createStore);

export default function configureStore(initialState) {
  return createFinalStore(reducer, initialState);
}
