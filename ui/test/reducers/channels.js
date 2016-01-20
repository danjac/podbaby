import immutable from 'immutable';
import { assert } from 'chai';
import { Actions } from '../../constants';
import channelsReducer from '../../reducers/channels';


describe('Channels', function () {
  it('Shows filtered channels', function () {
    const state = immutable.fromJS({
      channels: [
        {
          id: 100,
          title: 'test1',
        },
      ],
      filter: '',
    });

    const action = {
      type: Actions.FILTER_CHANNELS,
      payload: 'foo',
    };

    const newState = channelsReducer(state, action);

    assert.equal(newState.get('filter'), 'foo');
  });

  it('Shows all channels if filter is empty', function () {
    const state = immutable.fromJS({
      channels: [
        {
          id: 100,
          title: 'test1',
        },
      ],
      filter: null,
    });

    const action = {
      type: Actions.FILTER_CHANNELS,
      payload: '',
    };

    const newState = channelsReducer(state, action);

    assert.equal(newState.get('filter'), '');
  });
});
