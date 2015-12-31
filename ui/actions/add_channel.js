import * as api from '../api';
import * as alerts from './alerts';
import { Actions } from '../constants';
import { createAction } from './utils';

export const open = () => createAction(Actions.OPEN_ADD_CHANNEL_FORM);
export const close = () => createAction(Actions.CLOSE_ADD_CHANNEL_FORM);

export function add(url) {
    return dispatch => {
        api.addChannel(url)
        .then(result => {
            dispatch(alerts.success("New channel added: see your subscriptions"));
            dispatch(createAction(Actions.ADD_CHANNEL_SUCCESS, result.data));
        })
        .catch(error => {
            dispatch(createAction(Actions.ADD_CHANNEL_FAILURE, { error }));
        });
    };
}
