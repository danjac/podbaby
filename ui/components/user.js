import React, { PropTypes } from 'react';
import { connect } from 'react-redux';

import {
    ButtonInput,
    Input
} from 'react-bootstrap';

export class User extends React.Component {

  render() {
    return (
      <div>
        <h3>Change my email address</h3>
        <form className="form form-vertical">
            <Input ref="email"
                   type="email"
                   required
                   defaultValue={this.props.auth.email}  />
            <ButtonInput bsStyle="primary"
                         value="Save" />
        </form>
        <h3>Change my password</h3>
        <form className="form form-vertical">
            <Input ref="oldPassword"
                   type="password"
                   placeholder="Old password"
                   required />
            <Input ref="newPassword"
                   type="password"
                   placeholder="New password"
                   required />
            <ButtonInput bsStyle="primary"
                         value="Save" />
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
