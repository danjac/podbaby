import React, { PropTypes } from 'react';
import { connect } from 'react-redux';

export class Search extends React.Component {
  render() {
    return (
      <div>
        <h2>Searching for {this.props.query}</h2>

      </div>
    );
  }
}

const mapStateToProps = state => {
  const { q } = state.search;
  return {
    query: q
  };
};

export default connect(mapStateToProps)(Search);



