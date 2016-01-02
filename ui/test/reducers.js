import { assert } from 'chai';
import { Actions } from '../constants';
import channelsReducer from '../reducers/channels';
import alertsReducer from '../reducers/alerts';

describe('Channels', function() {

  it('Shows filtered channels', function() {

    const state = {
      requestedChannels: [
        {
          id: 100,
          title: 'test1'
        }
      ],
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
      payload: 'foo'
    };

    const newState = channelsReducer(state, action);

    assert.equal(newState.channels.length, 0);
    // requestedChannels should remain the same
    assert.equal(newState.requestedChannels.length, 1);
    assert.notEqual(newState.filter, null);

  });

  it('Shows all channels if filter is empty', function() {

    const state = {
      requestedChannels: [
        {
          id: 100,
          title: 'test1'
        }
      ],
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

    assert.equal(newState.channels.length, 1);
    assert.equal(newState.requestedChannels.length, 1);
    assert.equal(newState.filter, null);

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
