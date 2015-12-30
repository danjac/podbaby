import { Actions } from '../constants';
import { createAction } from './utils';

export const requestPodcasts = () => createAction(Actions.PODCASTS_REQUEST);
