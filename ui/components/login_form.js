import React, { PropTypes } from 'react';

import { reduxForm } from 'redux-form';


import Button from 'react-bootstrap';

import Icon from '../components/icon';
import { FormGroup } from '../components/form';


const validateLogin = values => {
  const errors = {};
  if (!values.identifier) {
    errors.identifier = 'Email or name required';
  }
  if (!values.password) {
    errors.password = 'Password is required';
  }
  return errors;
};

export class LoginForm extends React.Component {

  render() {
    const {
      fields: { identifier, password },
      handleSubmit,
      submitting,
    } = this.props;

    return (
      <form className="form-horizontal" onSubmit={handleSubmit}>

        <FormGroup field={identifier}>
          <input
            type="text"
            className="form-control"
            placeholder="Email address or name"
            {...identifier}
          />
        </FormGroup>

        <FormGroup field={password}>
          <input
            type="password"
            className="form-control"
            placeholder="Password"
            {...password}
          />
        </FormGroup>

          <Button
            bsStyle="primary"
            disabled={submitting}
            onClick={handleSubmit}
            className="form-control"
            type="submit"
          >
            <Icon icon="sign-in" /> Login
          </Button>
      </form>
    );
  }
}

LoginForm.propTypes = {
  fields: PropTypes.object.isRequired,
  handleSubmit: PropTypes.func.isRequired,
  submitting: PropTypes.bool.isRequired,
};

LoginForm = reduxForm({
  form: 'login',
  fields: ['identifier', 'password'],
  validate: validateLogin,
})(LoginForm);
