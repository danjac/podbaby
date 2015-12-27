import * as api from '../api';
import { Actions } from '../constants';
import { createAction } from './utils';

export function getChannel(id) {
    return dispatch => {
        api.getChannel(id)
        .then(result => {
            dispatch(createAction(Actions.GET_CHANNEL_SUCCESS, result.data));
        })
        .catch(error => {
            dispatch(createAction(Actions.GET_CHANNEL_FAILURE, { error }));
        });
    };
}