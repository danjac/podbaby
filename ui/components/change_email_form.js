import React, { PropTypes } from 'react';
import { reduxForm } from 'redux-form';
import validator from 'validator';

import { Button } from 'react-bootstrap';

import * as api from '../api';
import Icon from '../components/icon';
import { FormGroup } from '../components/form';

const validateChangeEmail = values => {
  return values.email && validator.isEmail(values.email) ? {} :
    { email: 'Please provide a valid email address' };
};

export default class ChangeEmailForm extends React.Component {

  constructor(props) {
    super(props);
    this.handleSubmit = this.props.handleSubmit(this.handleSubmit.bind(this));
  }

  handleSubmit(values) {
    return new Promise((resolve, reject) => {
      const { email } = values;
      api.changeEmail(values.email)
      .then(() => {
        this.props.onComplete(email);
        resolve();
      }, error => {
        reject(error.data);
      });
    });
  }

  render() {
    const {
      submitting,
      fields: { email },
    } = this.props;

    return (
      <form className="form form-vertical" onSubmit={this.handleSubmit}>
            <FormGroup field={email}>
            <input type="email" className="form-control" {...email} />
          </FormGroup>

          <Button
            bsStyle="primary"
            className="form-control"
            disabled={submitting}
            type="submit"
          >
            <Icon icon="save"/> Save new email address
          </Button>
      </form>
    );
  }
}

ChangeEmailForm.propTypes = {
  onComplete: PropTypes.func.isRequired,
  handleSubmit: PropTypes.func.isRequired,
  submitting: PropTypes.bool.isRequired,
  fields: PropTypes.object.isRequired,
};

ChangeEmailForm = reduxForm({
  form: 'change-email',
  fields: ['email'],
  validate: validateChangeEmail,
}, state => ({
  initialValues: state.auth,
}))(ChangeEmailForm);
