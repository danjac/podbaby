import React, { PropTypes } from 'react';

import DocumentTitle from 'react-document-title';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import { reduxForm } from 'redux-form';


import {
  Input,
  Button,
  ButtonGroup,
  Modal
} from 'react-bootstrap';

import * as actions from '../actions';
import Icon from '../components/icon';
import { FormGroup } from '../components/form';
import { getTitle } from './utils';

const validateRecoverPassword = values => {
  return values.identifier ? {} : { identifier: 'You must provide a name or email' };
};


export class RecoverPasswordModal extends React.Component {

  render() {

    const {
      fields: { identifier },
      handleSubmit,
      submitting,
      onSubmit,
      resetForm,
      show,
      onClose,
      container
    } = this.props;

    const handleOnClose = () => {
      resetForm();
      onClose();
    };

    const handleOnSubmit = () => {
      handleSubmit();
      resetForm();
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
            <form className="form" onSubmit={handleOnSubmit}>
              <FormGroup field={identifier}>
                  <input type="text" className="form-control" placeholder="Email address or name" {...identifier} />
              </FormGroup>
              <ButtonGroup>
                <Button bsStyle="primary"
                        disabled={submitting}
                        onClick={handleOnSubmit}
                        type="submit"><Icon icon="send" /> Send</Button>
                <Button bsStyle="default" onClick={handleOnClose}><Icon icon="remove" /> Cancel</Button>
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
  }

  handleLogin(values) {
    const { identifier, password } = values;
    this.actions.login(identifier, password);
  }

  handleRecoverPassword(values) {
    this.actions.recoverPassword(values.identifier);
  }

  handleOpenRecoverPasswordForm(event) {
    event.preventDefault();
    this.actions.openRecoverPasswordForm();
  }

  handleCloseRecoverPasswordForm() {
    this.actions.closeRecoverPasswordForm();
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
                              onSubmit={this.handleRecoverPassword.bind(this)}
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
