import _ from 'lodash';
import React, { PropTypes } from 'react';
import { connect } from 'react-redux';
import DocumentTitle from 'react-document-title';

import { ListGroup, ListGroupItem } from 'react-bootstrap';
import { getTitle } from './utils';

export class Categories extends React.Component {
  render() {
    const { categories } = this.props;
    const { createHref } = this.props.history;
    return (
      <DocumentTitle title={getTitle('Browse categories')}>
        <ListGroup>
          {categories.map(category => {
            return (
            <ListGroupItem
                key={category.id}
                href={createHref(`/categories/${category.id}/`)}
              >
                {_.capitalize(category.name)}
              </ListGroupItem>
            );
          })}
        </ListGroup>
      </DocumentTitle>
    );
  }
}

Categories.propTypes = {
  categories: PropTypes.array.isRequired,
  history: PropTypes.object.isRequired,
};

const mapStateToProps = state => {
  const { categories } = state.categories;
  return { categories };
};

export default connect(mapStateToProps)(Categories);
