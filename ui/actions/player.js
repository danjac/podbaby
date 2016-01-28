import _ from 'lodash';
import * as api from '../api';
import { Actions, Storage } from '../constants';
import { createAction } from './utils';

const removePlayerFromSession = () => {
  window.sessionStorage.removeItem(Storage.CURRENT_PODCAST);
};

const savePlayerToSession = player => {
  window.sessionStorage.setItem(Storage.CURRENT_PODCAST, JSON.stringify(player));
};

export function updateTime(currentTime) {
  return (dispatch, getState) => {
    const { player } = getState();
    dispatch(createAction(Actions.PLAYER_TIME_UPDATE, currentTime));
    savePlayerToSession(player);
  };
}

export function close() {
  removePlayerFromSession();
  return createAction(Actions.CLOSE_PLAYER);
}

export function togglePlayer(podcast) {
  if (podcast.isPlaying) {
    return close();
  }

  return (dispatch, getState) => {
    dispatch(createAction(Actions.CURRENTLY_PLAYING, podcast));
    const { player, auth } = getState();
    if (auth.isLoggedIn) {
      api.nowPlaying(podcast.id);
    }
    savePlayerToSession(player);
  };
}

function playPodcastWithBookmark(id, dispatch, state) {
  const onSuccess = result => {
    dispatch(createAction(Actions.UPDATE_PODCAST_CACHE, result));
    dispatch(createAction(Actions.BOOKMARKS_CURRENTLY_PLAYING, result.id));
    dispatch(togglePlayer(Object.assign({}, result, { isBookmarked: true })));
  };

  const { podcastCache } = state.podcasts;
  const cached = podcastCache[id];
  if (cached) {
    onSuccess(cached);
    return;
  }

  api.getPodcast(id)
  .then(result => {
    onSuccess(result.data);
  })
  .catch(error => {
    dispatch(createAction(Actions.GET_PODCAST_FAILURE, { error }));
  });
}


function playNextBookmarkedPodcast(pos) {
  return (dispatch, getState) => {
    const state = getState();
    const { bookmarks, playing } = state.bookmarks;

    let nextPlaying;

    if (playing) {
      let index = bookmarks.indexOf(playing) + pos;
      if (index < 0) {
        index = bookmarks.length - 1;
      }
      nextPlaying = bookmarks[index] || bookmarks[0];
    } else {
      nextPlaying = bookmarks[0];
    }

    if (nextPlaying) {
      playPodcastWithBookmark(nextPlaying, dispatch, state);
    }
  };
}

export function playRandom() {
  return (dispatch, getState) => {
    const state = getState();
    const { bookmarks } = state.bookmarks;
    if (!bookmarks) {
      return;
    }
    const bookmarkId = _.sample(bookmarks);
    playPodcastWithBookmark(bookmarkId, dispatch, state);
  };
}

export function playLast() {
  return playNextBookmarkedPodcast(-1);
}


export function playNext() {
  return playNextBookmarkedPodcast(1);
}

// reload player from session
export function reloadPlayer() {
  const data = window.sessionStorage.getItem(Storage.CURRENT_PODCAST);
  const player = data ? JSON.parse(data) : null;
  return createAction(Actions.RELOAD_PLAYER, player);
}
