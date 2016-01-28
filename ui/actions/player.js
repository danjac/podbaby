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

function playBookmarked(pos) {
  return (dispatch, getState) => {
    const state = getState();
    const { bookmarks } = state.bookmarks;
    const playing = state.player.bookmarkId;

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
      const onSuccess = podcast => {
        dispatch(createAction(Actions.UPDATE_PODCAST_CACHE, podcast));
        dispatch(createAction(Actions.BOOKMARKS_CURRENTLY_PLAYING, podcast.id));
        dispatch(togglePlayer(Object.assign({}, podcast, { isBookmarked: true })));
      };

      const { podcastCache } = state.podcasts;
      const cached = podcastCache[nextPlaying];
      if (cached) {
        onSuccess(cached);
        return;
      }

      api.getPodcast(nextPlaying)
      .then(result => {
        onSuccess(result.data);
      })
      .catch(error => {
        dispatch(createAction(Actions.GET_PODCAST_FAILURE, { error }));
      });
    }
  };
}

export function playLast() {
  return playBookmarked(-1);
}


export function playNext() {
  return playBookmarked(1);
}

// reload player from session
export function reloadPlayer() {
  const data = window.sessionStorage.getItem(Storage.CURRENT_PODCAST);
  const player = data ? JSON.parse(data) : null;
  return createAction(Actions.RELOAD_PLAYER, player);
}
