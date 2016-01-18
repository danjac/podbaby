import React, { PropTypes } from 'react';

export class Wrapper extends React.Component {
  render() {
    return (
      <div>{this.props.children}</div>
    );
  }
}

Wrapper.propTypes = {
  children: PropTypes.object.isRequired,
};
