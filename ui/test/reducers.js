import { assert } from 'chai';
import { Actions, Alerts } from '../constants';
import channelsReducer from '../reducers/channels';
import alertsReducer from '../reducers/alerts';

describe('Channels', function() {

  it('Shows filtered channels', function() {

    const state = {
      channels: [
        {
          id: 100,
          title: 'test1'
        }
      ],
      filter: ""
    };

    const action = {
      type: Actions.FILTER_CHANNELS,
      payload: 'foo'
    };

    const newState = channelsReducer(state, action);

    assert.equal(newState.filter, "foo");

  });

  it('Shows all channels if filter is empty', function() {

    const state = {
      channels: [
        {
          id: 100,
          title: 'test1'
        }
      ],
      filter: null
    };

    const action = {
      type: Actions.FILTER_CHANNELS,
      payload: ''
    };

    const newState = channelsReducer(state, action);

    assert.equal(newState.filter, "");

  });

});

describe("Add an alert", function() {

  it("Adds a new alert", function() {
      const state = [];

      const action = {
        type: Actions.ADD_ALERT,
        payload: {
          status: Alerts.SUCCESS,
          message: "it worked!"
        }
      };

      const newState = alertsReducer(state, action);
      assert.equal(newState.length, 1);
  });

});

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
