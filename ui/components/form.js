import React from 'react';
import { Input } from 'react-bootstrap';

export const FormGroup = props => {

  const { field } = props;

  return (
    <Input hasFeedback={field.touched}
           bsStyle={field.touched ? (
           field.error ? 'error': 'success' ) : undefined}>
      {props.children}
      {field.touched && field.error && <span className="help-block">{field.error}</span>}
    </Input>
  );

};


