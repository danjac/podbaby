import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

import {
  Input,
  Button
} from 'react-bootstrap';

import { auth } from '../actions';

export class Signup extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindActionCreators(auth, dispatch);
  }

  handleSubmit(event) {

    event.preventDefault();

    const name = this.refs.name.getInputDOMNode().value;
    const email = this.refs.email.getInputDOMNode().value;
    const password = this.refs.password.getInputDOMNode().value;

    this.actions.signup(name, email, password);

  }
  render() {

    const { createHref } = this.props.history;

    return (
      <div>
        <h1>Join PodBaby today.</h1>
        <hr />
        <form className="form-horizontal"
              onSubmit={this.handleSubmit.bind(this)}>
            <Input required
              type="text"
              ref="name"
              placeholder="Name" />
            <Input required
              type="email"
              ref="email"
              placeholder="Email address" />
            <Input required
              type="password"
              ref="password"
              placeholder="Password" />
            <Button
              bsStyle="primary"
              className="form-control"
              type="submit">Signup</Button>
        </form>
        <p><a href={createHref('/login/')}>Already a member? Log in here.</a></p>
      </div>

    );
  }
};

Signup.propTypes = {
  dispatch: PropTypes.func.isRequired
};

export default connect()(Signup);
