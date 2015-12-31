import React, { PropTypes } from 'react';
import { connect } from 'react-redux';

import {
    ButtonInput,
    Input
} from 'react-bootstrap';

import * as actions from '../actions';

export class User extends React.Component {

  handleSubmitEmail(event) {
    event.preventDefault();
    const email = this.refs.email.getValue();
    if (email) {
      this.props.dispatch(actions.users.changeEmail(email));
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
      this.props.dispatch(actions.users.changePassword(oldPassword, newPassword));
    }

  }

  render() {
    return (
      <div>
        <h3>Change my username</h3>
        <form className="form form-vertical">
            <Input ref="name"
                   type="name"
                   required
                   defaultValue={this.props.auth.name}  />
            <ButtonInput bsStyle="primary"
                         type="submit"
                         value="Save" />
        </form>
        <h3>Change my email address</h3>
        <form className="form form-vertical" onSubmit={this.handleSubmitEmail.bind(this)}>
            <Input ref="email"
                   type="email"
                   required
                   defaultValue={this.props.auth.email}  />
            <ButtonInput bsStyle="primary"
                         type="submit"
                         value="Save" />
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

            <ButtonInput bsStyle="primary"
                         type="submit"
                         value="Save" />
        </form>
        <hr />
        <form>
          <ButtonInput bsStyle="danger" className="form-control" value="Delete my account" />
        </form>
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
