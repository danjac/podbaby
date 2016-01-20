import React, { PropTypes } from 'react';
import _ from 'lodash';
import md5 from 'md5';


const Gravatar = props => {
  const size = props.size || 20;
  const code = md5(_.trim(props.email.toLowerCase()));
  const src = `https://gravatar.com/avatar/${code}?s=${size}`;
  return <img src={src} />;
};

Gravatar.propTypes = {
  email: PropTypes.string.isRequired,
  size: PropTypes.number,
};

export default Gravatar;
