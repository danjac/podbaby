import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import DocumentTitle from 'react-document-title';
import { Link } from 'react-router';

import {
  Input,
  Button
} from 'react-bootstrap';

import { auth } from '../actions';
import Icon from '../components/icon';
import { getTitle } from './utils';

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

    return (
    <DocumentTitle title={getTitle('Signup')}>
      <div>
        <h2>Join PodBaby today.</h2>
        <hr />
        <p className="lead">
          As a member you can subscribe to podcast feeds and keep track of your favorite episodes.
        </p>
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
              type="submit"><Icon icon="sign-in" /> Signup</Button>
        </form>
        <p><Link to='/login/'>Already a member? Log in here.</Link></p>
      </div>
    </DocumentTitle>

    );
  }
};

Signup.propTypes = {
  dispatch: PropTypes.func.isRequired
};

export default connect()(Signup);
