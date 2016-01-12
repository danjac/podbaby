import React from 'react';
import classnames from 'classnames';

export default function(props) {
  const classes = classnames('fa', 'fa-' + props.icon, {'fa-spin': props.spin});
  return <i className={classes} />
}
