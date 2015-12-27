import { Actions } from '../constants';

const initialState = [];

export default function(state=initialState, action) {
    switch(action.type) {
        case Actions.GET_CHANNELS_SUCCESS:
            return action.payload;
        case Actions.GET_CHANNELS_FAILURE:
            return initialState;
    }

    return state;
}