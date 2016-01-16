import React, { PropTypes } from 'react';

import DocumentTitle from 'react-document-title';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router';

import * as actions from '../actions';
import LoginForm from '../components/login_form';
import RecoverPasswordModal from '../components/recover_password';
import { getTitle } from './utils';


export class Login extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindActionCreators(actions.auth, dispatch);
    this.alerts = bindActionCreators(actions.alerts, dispatch);

    this.handleLoginComplete = this.handleLoginComplete.bind(this);
    this.handleOpenRecoverPasswordForm = this.handleOpenRecoverPasswordForm.bind(this);
    this.handleCloseRecoverPasswordForm = this.handleCloseRecoverPasswordForm.bind(this);
    this.handleRecoverPasswordComplete = this.handleRecoverPasswordComplete.bind(this);
  }

  handleLoginComplete(result) {
    this.actions.loginComplete(result);
  }

  handleOpenRecoverPasswordForm(event) {
    event.preventDefault();
    this.actions.openRecoverPasswordForm();
  }

  handleCloseRecoverPasswordForm() {
    this.actions.closeRecoverPasswordForm();
  }

  handleRecoverPasswordComplete() {
    this.actions.recoverPasswordComplete();
  }

  render() {
    return (
    <DocumentTitle title={getTitle('Login')}>
      <div>
        <h2>Sign into your PodBaby account.</h2>
        <hr />
        <LoginForm onComplete={this.handleLoginComplete} />
        <p>
          <a href="#" onClick={this.handleOpenRecoverPasswordForm}>Forgot your password?</a>
        </p><p>
          <Link to="/signup/">Not a member yet? Sign up today!</Link>
        </p>
        <RecoverPasswordModal
          container={this}
          show={this.props.auth.showRecoverPasswordForm}
          onComplete={this.handleRecoverPasswordComplete}
          onClose={this.handleCloseRecoverPasswordForm}
        />
      </div>
    </DocumentTitle>

    );
  }
}

Login.propTypes = {
  auth: PropTypes.object.isRequired,
  dispatch: PropTypes.func.isRequired,
};

const mapStateToProps = state => {
  return {
    auth: state.auth,
  };
};

export default connect(mapStateToProps)(Login);
