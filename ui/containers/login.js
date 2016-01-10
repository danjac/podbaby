import React, { PropTypes } from 'react';

import DocumentTitle from 'react-document-title';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import { reduxForm } from 'redux-form';


import {
  Input,
  Button,
  Alert,
  ButtonGroup,
  Modal
} from 'react-bootstrap';

import * as actions from '../actions';
import * as api from '../api';
import Icon from '../components/icon';
import { FormGroup } from '../components/form';
import { getTitle } from './utils';

const validateRecoverPassword = values => {
  return values.identifier ? {} : { identifier: 'You must provide a name or email' };
};


export class RecoverPasswordModal extends React.Component {

  handleSubmit(values) {

    const { identifier } = values;
    const { resetForm, onComplete } = this.props;

    return new Promise((resolve, reject) => {

      return api.recoverPassword(identifier)
      .then(result => {
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
      fields: { identifier },
      handleSubmit,
      resetForm,
      submitting,
      onSubmit,
      show,
      onClose,
      container,
    } = this.props;

    const handleClose = () => {
      resetForm();
      onClose();
    };

    return (
      <Modal show={show}
             aria-labelledby="recover-password-modal-title"
             container={container}
             onHide={onClose}>
        <Modal.Header closeButton>
          <Modal.Title id="recover-password-modal-title">Recover password</Modal.Title>
        </Modal.Header>
        <Modal.Body>
              <p>We'll send you a new random password so you can log back in again.</p>
            <form className="form" onSubmit={handleSubmit(this.handleSubmit.bind(this))}>
              <FormGroup field={identifier}>
                  <input type="text" className="form-control" placeholder="Email address or name" {...identifier} />
              </FormGroup>
              <ButtonGroup>
                <Button bsStyle="primary"
                        disabled={submitting}
                        type="submit"><Icon icon="send" /> Send</Button>
                <Button bsStyle="default" onClick={handleClose}><Icon icon="remove" /> Cancel</Button>
              </ButtonGroup>
            </form>
        </Modal.Body>
      </Modal>
    );
  }

}

RecoverPasswordModal = reduxForm({
  form: 'recover-password',
  fields: ['identifier'],
  validate: validateRecoverPassword
})(RecoverPasswordModal);

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
      submitting
    } = this.props;

    return (
      <form className="form-horizontal" onSubmit={handleSubmit}>

        <FormGroup field={identifier}>
          <input type="text" className="form-control" placeholder="Email address or name" {...identifier} />
        </FormGroup>

        <FormGroup field={password}>
          <input type="password" className="form-control" placeholder="Password" {...password} />
        </FormGroup>

          <Button
            bsStyle="primary"
            disabled={submitting}
            onClick={handleSubmit}
            className="form-control"
            type="submit"><Icon icon="sign-in" /> Login</Button>
      </form>
    );
  }
}

LoginForm = reduxForm({
  form: 'login',
  fields: ['identifier', 'password'],
  validate: validateLogin
})(LoginForm);

export class Login extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindActionCreators(actions.auth, dispatch);
    this.alerts = bindActionCreators(actions.alerts, dispatch);
  }

  handleLogin(values) {
    const { identifier, password } = values;
    return new Promise((resolve, reject) => {
      return api.login(identifier, password)
      .then(result => {
        this.actions.loginComplete(result.data);
        resolve();
      }, error => {
        this.alerts.warning("Sorry, you were unable to log in");
        reject(error.data);
      });
    });
  }

  handleOpenRecoverPasswordForm(event) {
    event.preventDefault();
    this.actions.openRecoverPasswordForm();
  }

  handleCloseRecoverPasswordForm() {
    this.actions.closeRecoverPasswordForm();
  }

  handleRecoverPassword() {
    this.actions.recoverPasswordComplete();
  }

  render() {

    return (
    <DocumentTitle title={getTitle("Login")}>
      <div>
        <h2>Sign into your PodBaby account.</h2>
        <hr />
        <LoginForm onSubmit={this.handleLogin.bind(this)} />
        <p>
          <a href="#" onClick={this.handleOpenRecoverPasswordForm.bind(this)}>Forgot your password?</a>
        </p><p>
          <Link to="/signup/">Not a member yet? Sign up today!</Link>
        </p>
        <RecoverPasswordModal show={this.props.auth.showRecoverPasswordForm}
                              container={this}
                              onComplete={this.handleRecoverPassword.bind(this)}
                              onClose={this.handleCloseRecoverPasswordForm.bind(this)} />
      </div>
    </DocumentTitle>

    );
  }
};

Login.propTypes = {
  auth: PropTypes.object.isRequired,
  dispatch: PropTypes.func.isRequired
};

const mapStateToProps = state => {
  return {
    auth: state.auth
  };
}

export default connect(mapStateToProps)(Login);
