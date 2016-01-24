import React, { PropTypes } from 'react';
import _ from 'lodash';
import md5 from 'md5';


const Gravatar = ({ email, size }) => {
  const code = md5(_.trim(email.toLowerCase()));
  const src = `https://gravatar.com/avatar/${code}?s=${size || 20}`;
  return <img src={src} />;
};

Gravatar.propTypes = {
  email: PropTypes.string.isRequired,
  size: PropTypes.number,
};

export default Gravatar;
