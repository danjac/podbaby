import { Actions } from '../constants';
import { createAction } from './utils';

export const unloadPodcasts = () => createAction(Actions.UNLOAD_PODCASTS);
