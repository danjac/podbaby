import { assert } from 'chai';
import { Actions } from '../constants';
import alertsReducer from '../reducers/alerts';

describe('Dismiss an alert', function() {

  it('Removes a messaage if ID found', function() {
    const state =  [
      {
        id: 1000,
        status: "info",
        message: "testing"
      }
    ];
    const action = {
      type: Actions.DISMISS_ALERT,
      payload: 1000
    }
    const newState = alertsReducer(state, action);
    assert.equal(newState.length, 0)

  });

  it('Does nothing if no matching ID', function() {
    const state =  [
      {
        id: 1000,
        status: "info",
        message: "testing"
      }
    ];
    const action = {
      type: Actions.DISMISS_ALERT,
      payload: 1001
    }
    const newState = alertsReducer(state, action);
    assert.equal(newState.length, 1)

  });
});
