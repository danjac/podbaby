import React, { PropTypes } from 'react';
import { reduxForm } from 'redux-form';
import validator from 'validator';

import { Button } from 'react-bootstrap';

import * as api from '../api';
import Icon from '../components/icon';
import { FormGroup } from '../components/form';


const validateChangePassword = values => {
  const errors = {};

  if (!values.oldPassword) {
    errors.oldPassword = 'Please provide your old password';
  }
  if (!validator.isLength(values.newPassword, 6)) {
    errors.newPassword = 'Your new password must be at least 6 characters in length';
  }
  return errors;
};


export class ChangePasswordForm extends React.Component {

  handleSubmit(values) {
    const { onComplete, resetForm } = this.props;
    return new Promise((resolve, reject) => {
      api.changePassword(values.oldPassword, values.newPassword)
      .then(() => {
        onComplete();
        resetForm();
        resolve();
      }, error => {
        reject(error.data);
      });
    });
  }

  render() {
    const {
      handleSubmit,
      fields: { oldPassword, newPassword },
    } = this.props;

    return (
        <form className="form form-vertical" onSubmit={handleSubmit(this.handleSubmit.bind(this))}>

          <FormGroup field={oldPassword}>
            <input
              type="password"
              className="form-control"
              placeholder="Your old password"
              {...oldPassword}
            />
          </FormGroup>

          <FormGroup field={newPassword}>
            <input
              type="password"
              className="form-control"
              placeholder="Your old password"
              {...newPassword}
            />
          </FormGroup>

          <Button
            bsStyle="primary"
            className="form-control"
            type="submit"
          ><Icon icon="save" /> Save new password
          </Button>
        </form>

    );
  }
}

ChangePasswordForm.propTypes = {
  onComplete: PropTypes.func.isRequired,
  handleSubmit: PropTypes.func.isRequired,
  resetForm: PropTypes.func.isRequired,
  fields: PropTypes.object.isRequired,
  submitting: PropTypes.bool.isRequired,
};

ChangePasswordForm = reduxForm({
  form: 'change-password',
  fields: ['oldPassword', 'newPassword'],
  validate: validateChangePassword,
})(ChangePasswordForm);
