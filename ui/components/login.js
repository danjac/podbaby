import React, { PropTypes } from 'react';

import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router';


import {
  Input,
  Button,
  ButtonGroup,
  Modal,
  Glyphicon
} from 'react-bootstrap';

import * as actions from '../actions';

export class RecoverPasswordModal extends React.Component {

  handleSubmit(event){
    event.preventDefault();
    const node = this.refs.identifier.getInputDOMNode();
    if (node.value) {
      this.props.onSubmit(node.value);
      node.value = "";
    }
  }

  render() {
    const { show, onClose, container } = this.props;
    return (
      <Modal show={show}
             aria-labelledby="recover-password-modal-title"
             container={container}
             onHide={onClose}>
        <Modal.Header closeButton>
          <Modal.Title id="recover-password-modal-title">Recover password</Modal.Title>
        </Modal.Header>
        <Modal.Body>
            <form className="form" onSubmit={this.handleSubmit.bind(this)}>
              <Input required type="text" placeholder="Your username or email address" ref="identifier" />
              <ButtonGroup>
                <Button bsStyle="primary" type="submit"><Glyphicon glyph="plus" /> Submit</Button>
                <Button bsStyle="default" onClick={onClose}><Glyphicon glyph="remove" /> Cancel</Button>
              </ButtonGroup>
            </form>
        </Modal.Body>
      </Modal>
    );
  }

}

export class Login extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindActionCreators(actions.auth, dispatch);
  }

  handleLogin(event) {
    event.preventDefault();
    const identifier = this.refs.identifier.getInputDOMNode().value;
    const password = this.refs.password.getInputDOMNode().value;
    if (identifier && password) {
      this.actions.login(identifier, password);
    }
  }

  handleOpenRecoverPasswordForm(event) {
    event.preventDefault();
    this.actions.openRecoverPasswordForm();
  }

  handleRecoverPassword(identifier) {
    this.actions.recoverPassword(identifier);
  }

  handleCloseRecoverPasswordForm() {
    this.actions.closeRecoverPasswordForm();
  }

  render() {

    return (
      <div>
        <h1>Sign into your PodBaby account.</h1>
        <hr />
        <form className="form-horizontal" onSubmit={this.handleLogin.bind(this)}>
            <Input required
              type="text"
              ref="identifier"
              placeholder="Email or username" />
            <Input required
              type="password"
              ref="password"
              placeholder="Password" />
            <Button
              bsStyle="primary"
              className="form-control"
              type="submit">Login</Button>
        </form>
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
