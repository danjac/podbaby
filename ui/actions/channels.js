import * as api from '../api';
import { Actions } from '../constants';
import { createAction } from './utils';

export function getChannels() {
    return dispatch => {
        api.getChannels()
        .then(result => {
            dispatch(createAction(Actions.GET_CHANNELS_SUCCESS, result.data));
        })
        .catch(error => {
            dispatch(createAction(Actions.GET_CHANNELS_FAILURE, { error }));
        });
    };
}