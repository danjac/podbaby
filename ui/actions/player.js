import { Actions } from '../constants';
import { createAction } from './utils';

export const setPodcast = podcast => createAction(Actions.CURRENTLY_PLAYING, podcast);
