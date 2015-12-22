import React, { PropTypes } from 'react';

import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

import {
  Input,
  Button
} from 'react-bootstrap';

import * as actions from '../actions';

export class Login extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindActionCreators(actions.auth, dispatch);
  }

  login(event) {
    event.preventDefault();
    const identifier = this.refs.identifier.getInputDOMNode();
    const password = this.refs.password.getInputDOMNode();
    this.actions.login(identifier, password);
  }

  render() {

    return (
      <div>
        <form className="form-horizontal" onSubmit={this.login.bind(this)}>
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
      </div>

    );
  }
};

Login.propTypes = {
  dispatch: PropTypes.func.isRequired
};

export default connect()(Login);
