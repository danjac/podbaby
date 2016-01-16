import React, { PropTypes } from 'react';
import { reduxForm } from 'redux-form';
import { Button, Alert } from 'react-bootstrap';

import * as api from '../api';
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

  constructor(props) {
    super(props);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleSubmit(values) {
    const { identifier, password } = values;
    return new Promise((resolve, reject) => {
      return api.login(identifier, password)
      .then(result => {
        this.props.onComplete(result.data);
        resolve();
      }, () => {
        reject({ _error: 'Sorry, unable to log  you in' });
      });
    });
  }

  render() {
    const {
      fields: { identifier, password },
      handleSubmit,
      submitting,
      error,
    } = this.props;

    return (
      <form className="form-horizontal" onSubmit={handleSubmit(this.handleSubmit)}>
        {error ?
          <Alert
            bsStyle="danger"
            className="text-center"
          >{error}</Alert> : ''}
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
          onClick={handleSubmit(this.handleSubmit)}
          className="form-control"
          type="submit"
        ><Icon icon="sign-in" /> Login
        </Button>
      </form>
    );
  }
}

LoginForm.propTypes = {
  error: PropTypes.string,
  fields: PropTypes.object.isRequired,
  handleSubmit: PropTypes.func.isRequired,
  onComplete: PropTypes.func.isRequired,
  submitting: PropTypes.bool.isRequired,
};

export default reduxForm({
  form: 'login',
  fields: ['identifier', 'password'],
  validate: validateLogin,
})(LoginForm);
