import React, { PropTypes } from 'react';

const PageHeader = ({ header }) => {
  return (
    <div className="page-header">
      <h3>{header}</h3>
    </div>
  );
};

PageHeader.propTypes = {
  header: PropTypes.any.isRequired,
};

export default PageHeader;
