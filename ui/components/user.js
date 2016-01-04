import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

import {
    Input,
    Button
} from 'react-bootstrap';

import * as actions from '../actions';
import Icon from './icon';

export class User extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindActionCreators(actions.users, dispatch);
  }

  handleSubmitEmail(event) {
    event.preventDefault();
    const email = this.refs.email.getValue();
    if (email) {
      this.actions.changeEmail(email);
    }

  }

  handleSubmitPassword(event) {
    event.preventDefault();

    const oldPasswordNode = this.refs.oldPassword.getInputDOMNode();
    const newPasswordNode = this.refs.newPassword.getInputDOMNode();

    const oldPassword = oldPasswordNode.value;
    const newPassword = newPasswordNode.value;

    if (oldPassword && newPassword){
      oldPasswordNode.value = "";
      newPasswordNode.value = "";
      this.actions.changePassword(oldPassword, newPassword);
    }

  }

  handleDelete(event) {
    event.preventDefault();
    if (window.confirm("Are you sure you want to delete this account? You will lose all your subscriptions and bookmarks!!!")) {
      this.actions.deleteAccount();
    }
  }

  render() {
    return (
      <div>
        <h3>Change my email address</h3>
        <form className="form form-vertical" onSubmit={this.handleSubmitEmail.bind(this)}>
            <Input ref="email"
                   type="email"
                   required
                   defaultValue={this.props.auth.email}  />
                 <Button bsStyle="primary" className="form-control" type="submit"><Icon icon="save" /> Save new email address</Button>
        </form>
        <h3>Change my password</h3>
        <form className="form form-vertical" onSubmit={this.handleSubmitPassword.bind(this)}>

            <Input ref="oldPassword"
                   type="password"
                   placeholder="Old password"
                   required />

            <Input ref="newPassword"
                   type="password"
                   placeholder="New password"
                   required />

                 <Button bsStyle="primary" className="form-control" type="submit"><Icon icon="save" /> Save new password</Button>
        </form>
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
