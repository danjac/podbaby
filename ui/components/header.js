import React, { PropTypes } from 'react';

const PageHeader = props => {
  return (
    <div className="page-header">
      <h2>{props.header}</h2>
    </div>
  );
};

PageHeader.propTypes = {
  header: PropTypes.any.isRequired,
};

export default PageHeader;
