import React from 'react';
import classnames from 'classnames';

export default function(props) {
  const classes = classnames('fa', 'fa-' + props.icon);
  return <i className={classes} />
}
