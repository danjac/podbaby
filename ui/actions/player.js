import * as api from '../api';
import { Actions, Storage } from '../constants';
import { createAction } from './utils';

const removePlayerFromSession = () => {
    window.sessionStorage.removeItem(Storage.CURRENT_PODCAST);
};

const savePlayerToSession = player => {
    window.sessionStorage.setItem(Storage.CURRENT_PODCAST, JSON.stringify(player));
};

export function updateTime(player, currentTime) {
  savePlayerToSession(player);
  return createAction(Actions.PLAYER_TIME_UPDATE, currentTime);
}

export function setPodcast(player, podcast) {
  if (podcast) {
    api.nowPlaying(podcast.id);
    savePlayerToSession(player);
  } else {
    removePlayerFromSession();
  }
  return createAction(Actions.CURRENTLY_PLAYING, podcast);
}

export function close(player, podcast) {
  removePlayerFromSession();
  return createAction(Actions.CLOSE_PLAYER);
}

// reload player from session
export function reloadPlayer() {
  const data = window.sessionStorage.getItem(Storage.CURRENT_PODCAST);
  const player = data ? JSON.parse(data) : null;
  return createAction(Actions.RELOAD_PLAYER, player);
}

