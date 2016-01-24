import React, { PropTypes } from 'react';

const PageHeader = ({ header }) => {
  return (
    <div className="page-header">
      <h2>{header}</h2>
    </div>
  );
};

PageHeader.propTypes = {
  header: PropTypes.any.isRequired,
};

export default PageHeader;
