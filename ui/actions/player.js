import { Actions } from '../constants';
import { Storage } from '../constants';
import { createAction } from './utils';

export function setPodcast(podcast) {
  if (podcast) {
    window.sessionStorage.setItem(Storage.CURRENT_PODCAST, JSON.stringify(podcast));
  } else {
    window.sessionStorage.removeItem(Storage.CURRENT_PODCAST);
  }
  return createAction(Actions.CURRENTLY_PLAYING, podcast);
}

// reload player from session
export function reloadPlayer() {
  const data = window.sessionStorage.getItem(Storage.CURRENT_PODCAST);
  const podcast = data ? JSON.parse(data) : null;
  return createAction(Actions.CURRENTLY_PLAYING, podcast);
}

