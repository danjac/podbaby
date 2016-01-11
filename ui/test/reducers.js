import { assert } from 'chai';
import { Actions, Alerts } from '../constants';
import channelsReducer from '../reducers/channels';
import alertsReducer from '../reducers/alerts';
import playerReducer from '../reducers/player';

describe("Player", function () {

  const initialState = {
    podcast: null,
    isPlaying: false,
    currentTime: 0
  };

  it('Sets the currently playing podcast', function() {

    // time should be reset to 0
    const state = Object.assign({}, initialState, { currentTime: 30 });

    const podcast = {
      id: 100,
      title: "A podcast"
    };

    const action = {
      type: Actions.CURRENTLY_PLAYING,
      payload: podcast
    };

    const newState = playerReducer(state, action);

    assert.equal(newState.podcast.id, 100);
    assert.equal(newState.currentTime, 0);
    assert.equal(newState.isPlaying, true);

  });

  it('Bookmarks the podcast', function() {
  });

  it('Bookmarks the podcast if player empty', function() {
  });

  it('Removes the bookmark from the podcast', function() {
  });

  it('Updates the current play time', function() {
  });

  it('Reloads the player if none currently playing', function() {
  });

  it('Reloads the player if podcast currently playing', function() {
  });

  it('Closes the player', function() {
  });
});

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

describe("Alerts", function() {

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
