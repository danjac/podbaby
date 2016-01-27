import { assert } from 'chai';
import { Actions } from '../../constants';
import podcastsReducer from '../../reducers/podcasts';

describe('Podcasts', function () {
  it('Hides podcast detail', function () {
    const state = {
      showDetail: [4],
    };

    const action = {
      type: Actions.HIDE_PODCAST_DETAIL,
      payload: 4,
    };

    const newState = podcastsReducer(state, action);
    assert.equal(newState.showDetail.length, 0);
  });

  it('Shows podcast detail', function () {
    const state = {
      showDetail: [],
    };

    const action = {
      type: Actions.SHOW_PODCAST_DETAIL,
      payload: 4,
    };

    const newState = podcastsReducer(state, action);
    assert.sameMembers(newState.showDetail, [4]);
  });

  it('Fails to return list of podcasts', function () {
    const state = {
      podcasts: [],
      showDetail: [],
      isLoading: true,
      page: {
        numPages: 0,
        numRows: 0,
        page: 1,
      },
    };

    const action = {
      type: Actions.LATEST_PODCASTS_FAILURE,
    };

    const newState = podcastsReducer(state, action);

    assert.notOk(newState.isLoading);
  });

  it('Returns list of podcasts', function () {
    const payload = {
      podcasts: [
        {
          id: 1000,
          title: 'testing',
        },
      ],
      page: {
        page: 1,
        numPages: 1,
        numRows: 1,
      },
    };

    const state = {
      podcasts: [],
      showDetail: [],
      isLoading: true,
      page: {
        numPages: 0,
        numRows: 0,
        page: 1,
      },
    };

    const action = {
      type: Actions.LATEST_PODCASTS_SUCCESS,
      payload,
    };

    const newState = podcastsReducer(state, action);

    assert.equal(newState.page, payload.page);
    assert.equal(newState.podcasts, payload.podcasts);
    assert.notOk(newState.isLoading);
  });
});
