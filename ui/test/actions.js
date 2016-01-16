import { assert } from 'chai';

import { Actions, Alerts } from '../constants';
import * as actions from '../actions';


describe('Alert actions', function () {
  it('Should fire a successful alert', function () {
    const result = actions.alerts.success('OK!!!');
    assert.equal(result.type, Actions.ADD_ALERT);
    const { message, status, id } = result.payload;
    assert.equal(status, Alerts.SUCCESS);
    assert.equal(message, 'OK!!!');
    assert(id, 'There should be an ID');
  });
});
