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

import * as api from '../api';
import * as actions from '../actions';
import Icon from '../components/icon';
import { FormGroup } from '../components/form';
import { getTitle } from './utils';

const validateChangeEmail = values => {
  return values.email && validator.isEmail(values.email) ? {} :
    { email: 'Please provide a valid email address' };
};

export class ChangeEmailForm extends React.Component {

  handleSubmit(values) {
    return new Promise((resolve, reject) => {
      const { email } = values;
      api.changeEmail(values.email)
      .then(result => {
        this.props.onComplete(email);
        resolve();
      }, error => {
        reject(error.data);
      });
    });
  }

  render() {
    const {
      handleSubmit,
      submitting,
      fields: { email }
    } = this.props;

    return (
      <form className="form form-vertical"
            onSubmit={handleSubmit(this.handleSubmit.bind(this))}>
            <FormGroup field={email}>
            <input type="email" className="form-control" {...email} />
          </FormGroup>

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
  validate: validateChangeEmail,
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
      submitting
    } = this.props;

    return (
        <form className="form form-vertical" onSubmit={handleSubmit(this.handleSubmit.bind(this))}>

          <FormGroup field={oldPassword}>
                <input type="password" className="form-control" placeholder="Your old password" {...oldPassword} />
          </FormGroup>

          <FormGroup field={newPassword}>
                <input type="password" className="form-control" placeholder="Your old password" {...newPassword} />
          </FormGroup>


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


export const DeleteAccountForm = props => {
  return (
      <div>
        <Button bsStyle="danger"
                className="form-control"
                onClick={props.onDelete}><Icon icon="trash" /> Delete my account</Button>
        <p className="text-center">
          <b>This will completely and irreversibly delete your account, including all your subscriptions and bookmarks.</b>
        </p>
      </div>
  );
};

export class User extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindActionCreators(actions.users, dispatch);
  }

  handleChangeEmailComplete(email) {
    this.actions.changeEmailComplete(email);
  }

  handleChangePasswordComplete() {
    this.actions.changePasswordComplete();
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
        <ChangeEmailForm onComplete={this.handleChangeEmailComplete.bind(this)} />
        <h3>Change my password</h3>
        <ChangePasswordForm onComplete={this.handleChangePasswordComplete.bind(this)} />
        <hr />
        <DeleteAccountForm onDelete={this.handleDelete.bind(this)} />
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
