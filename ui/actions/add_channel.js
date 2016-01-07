import _ from 'lodash';
import * as api from '../api';
import * as alerts from './alerts';
import { pushPath } from 'redux-simple-router';
import { Actions } from '../constants';
import { createAction } from './utils';

export const open = () => createAction(Actions.OPEN_ADD_CHANNEL_FORM);
export const close = () => createAction(Actions.CLOSE_ADD_CHANNEL_FORM);

export function add(url) {
    return dispatch => {
        dispatch(createAction(Actions.ADD_CHANNEL_REQUEST));
        api.addChannel(url)
        .then(result => {
            dispatch(alerts.success("New channel added"));
            dispatch(createAction(Actions.ADD_CHANNEL_SUCCESS, result.data));
            dispatch(pushPath(`/channel/${result.data.id}/`));
        })
        .catch(error => {
            dispatch(createAction(Actions.ADD_CHANNEL_FAILURE, { error }));
        });
    };
}
