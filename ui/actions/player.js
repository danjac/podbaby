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

  const payload = podcast.toJS();

  return (dispatch, getState) => {
    dispatch(createAction(Actions.CURRENTLY_PLAYING, payload));
    const { player, auth } = getState();
    if (auth.isLoggedIn) {
      api.nowPlaying(payload.id);
    }
    savePlayerToSession(player);
  };
}

// reload player from session
export function reloadPlayer() {
  const data = window.sessionStorage.getItem(Storage.CURRENT_PODCAST);
  const player = data ? JSON.parse(data) : null;
  return createAction(Actions.RELOAD_PLAYER, player);
}
