import _ from 'lodash';

export const createAction = (type, payload) => _.merge({ type }, { payload });
