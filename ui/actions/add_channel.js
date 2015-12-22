import { Actions } from '../constants';
import { createAction } from './utils';

export const open = () => createAction(Actions.OPEN_ADD_CHANNEL_FORM);

export const close = () => createAction(Actions.CLOSE_ADD_CHANNEL_FORM);
