import { Actions } from '../constants';
import { createAction } from './utils';

export function loadCategories(categories) {
  return createAction(Actions.CATEGORIES_LOADED, categories);
}
