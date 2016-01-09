import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import DocumentTitle from 'react-document-title';
import { connect } from 'react-redux';
import { reduxForm } from 'redux-form';
import validator from 'validator';

import {
    Input,
    Button
} from 'react-bootstrap';

import * as actions from '../actions';
import Icon from '../components/icon';
import { getTitle } from './utils';

const validateChangeEmail = values => {
  return values.email && validator.isEmail(values.email) ? {} :
    { email: 'Please provide a valid email address' };
};

export class ChangeEmailForm extends React.Component {

  render() {

    const {
      handleSubmit,
      submitting,
      fields: { email }
    } = this.props;

    return (
      <form className="form form-vertical"
            onSubmit={handleSubmit}>
          <Input hasFeedback={email.touched} bsStyle={email.touched ? ( email.error ? 'error': 'success' ) : undefined}>
                <input type="email" className="form-control" placeholder="Your new email address" {...email} />
                {email.touched && email.error && <div className="help-block">{email.error}</div>}
            </Input>
          <Button bsStyle="primary"
                  className="form-control"
                  disabled={submitting}
                  type="submit"><Icon icon="save" /> Save new email address</Button>
      </form>
    );

  }
}

ChangeEmailForm = reduxForm({
  form: 'change-email',
  fields: ['email'],
  validate: validateChangeEmail
}, state => ({
  initialValues: state.auth
}))(ChangeEmailForm);


const validateChangePassword = values => {
  const errors = {};

  if (!values.oldPassword) {
    errors.oldPassword = "Please provide your old password";
  }
  if (!validator.isLength(values.newPassword, 6)) {
    errors.newPassword = "Your new password must be at least 6 characters in length";
  }
  return errors;
};


export class ChangePasswordForm extends React.Component {
  render() {
    const {
      handleSubmit,
      fields: { oldPassword, newPassword },
      submitting,
      resetForm
    } = this.props;

    const onSubmit = () => {
      handleSubmit();
      resetForm();
    };

    return (
        <form className="form form-vertical" onSubmit={onSubmit}>
            <Input hasFeedback={oldPassword.touched} bsStyle={oldPassword.touched ? ( oldPassword.error ? 'error': 'success' ) : undefined}>
                <input type="password" className="form-control" placeholder="Your old password" {...oldPassword} />
                {oldPassword.touched && oldPassword.error && <div className="help-block">{oldPassword.error}</div>}
              </Input>
            <Input hasFeedback={newPassword.touched} bsStyle={newPassword.touched ? ( newPassword.error ? 'error': 'success' ) : undefined}>
                <input type="password" className="form-control" placeholder="Your new password" {...newPassword} />
                {newPassword.touched && newPassword.error && <div className="help-block">{newPassword.error}</div>}
              </Input>

            <Button bsStyle="primary" className="form-control" type="submit"><Icon icon="save" /> Save new password</Button>
        </form>

    );
  }
}

ChangePasswordForm = reduxForm({
  form: 'change-password',
  fields: ['oldPassword', 'newPassword'],
  validate: validateChangePassword
})(ChangePasswordForm);

export class User extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindActionCreators(actions.users, dispatch);
  }

  handleSubmitEmail(values) {
    this.actions.changeEmail(values.email);
  }

  handleSubmitPassword(values) {
    const { oldPassword, newPassword } = values;
    this.actions.changePassword(oldPassword, newPassword);
  }

  handleDelete(event) {
    event.preventDefault();
    if (window.confirm("Are you sure you want to delete this account? You will lose all your subscriptions and bookmarks!!!")) {
      this.actions.deleteAccount();
    }
  }

  render() {
    return (
    <DocumentTitle title={getTitle('My settings')}>
      <div>
        <h3>Change my email address</h3>
        <ChangeEmailForm onSubmit={this.handleSubmitEmail.bind(this)} />
        <h3>Change my password</h3>
        <ChangePasswordForm onSubmit={this.handleSubmitPassword.bind(this)} />
        <hr />
        <div>
          <Button bsStyle="danger"
                  className="form-control"
                  onClick={this.handleDelete.bind(this)}><Icon icon="trash" /> Delete my account</Button>
          <p className="text-center">
            <b>This will completely and irreversibly delete your account, including all your subscriptions and bookmarks.</b>
          </p>
        </div>
      </div>
    </DocumentTitle>
    );
  }

}

User.propTypes = {
  auth: PropTypes.object.isRequired,
  dispatch: PropTypes.func.isRequired
}

const mapStateToProps = state => {
    return {
      auth: state.auth
    };
};

export default connect(mapStateToProps)(User);
