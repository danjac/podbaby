import { bindActionCreators } from 'redux';

export const createAction = (type, payload) => {
  return { type, payload };
};

export const bindAllActionCreators = (actionCreators, dispatch) => {
  return Object.keys(actionCreators).reduce((result, key) => {
    return Object.assign({}, result, {
      [key]: bindActionCreators(actionCreators[key], dispatch),
    });
  }, {});
};
