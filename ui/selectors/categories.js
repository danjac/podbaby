import _ from 'lodash';
import { createSelector } from 'reselect';

const categoryMapSelector = state => state.categories.categoryMap;
const categoryPreSelector = state => state.categories.category;

const assignRelations = (c, categoryMap) => {
  const categories = _.values(categoryMap);
  return Object.assign({}, c, {
    children: categories.filter(child => child.parentId.Int64 === c.id),
    parent: _.find(categories, parent => parent.id === c.parentId.Int64),
  });
};

export const categoriesSelector = createSelector(
  [categoryMapSelector],
  (categoryMap) => {
    return _.sortBy(_.values(categoryMap).map(c => {
      return assignRelations(c, categoryMap);
    }), 'name');
  }
);

export const parentCategoriesSelector = createSelector(
  [categoriesSelector],
  (categories) => {
    return categories.filter(c => !c.parent);
  }
);

export const categorySelector = createSelector(
  [categoryMapSelector, categoryPreSelector],
  (categoryMap, category) => {
    return assignRelations(category, categoryMap);
  }
);
