import immutable from 'immutable';
import { assert } from 'chai';
import { Actions, Alerts } from '../../constants';
import alertsReducer from '../../reducers/alerts';


describe('Alerts', function () {
  it('Adds a new alert', function () {
    const state = immutable.List();

    const action = {
      type: Actions.ADD_ALERT,
      payload: {
        status: Alerts.SUCCESS,
        message: 'it worked!',
      },
    };

    const newState = alertsReducer(state, action);
    assert.equal(newState.size, 1);
  });

  it('Removes a message if ID found', function () {
    const state = immutable.fromJS([
      {
        id: 1000,
        status: 'info',
        message: 'testing',
      },
    ]);
    const action = {
      type: Actions.DISMISS_ALERT,
      payload: 1000,
    };
    const newState = alertsReducer(state, action);
    assert.equal(newState.size, 0);
  });

  it('Does nothing if no matching ID', function () {
    const state = immutable.fromJS([
      {
        id: 1000,
        status: 'info',
        message: 'testing',
      },
    ]);
    const action = {
      type: Actions.DISMISS_ALERT,
      payload: 1001,
    };
    const newState = alertsReducer(state, action);
    assert.equal(newState.size, 1);
  });
});
