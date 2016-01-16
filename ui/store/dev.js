import { createStore, compose } from 'redux';

import DevTools from '../containers/devtools';
import middleware from '../middleware';
import reducer from '../reducers';


const createFinalStore = compose(
  middleware,
  DevTools.instrument()
)(createStore);

export default function configureStore(initialState) {
  const store = createFinalStore(reducer, initialState);
  if (module.hot) {
    module.hot.accept('../reducers', () =>
      store.replaceReducer(require('../reducers'))
    );
  }
  return store;
}
