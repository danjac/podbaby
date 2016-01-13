import { createStore, compose } from 'redux';

import middleware from '../middleware';
import reducer from '../reducers';

const createFinalStore = compose(middleware)(createStore);

export default function configureStore(initialState) {
  return createFinalStore(reducer, initialState);
}
