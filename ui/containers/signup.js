import _ from 'lodash';
import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import DocumentTitle from 'react-document-title';
import { Link } from 'react-router';
import { reduxForm } from 'redux-form';
import validator from 'validator';

import {
  Input,
  Button
} from 'react-bootstrap';

import * as api from '../api';
import { auth } from '../actions';
import Icon from '../components/icon';
import { getTitle } from './utils';

const validate = values => {

  const { name, email, password } = values;
  const errors = {};

  if (!name || name.length < 3 || name.length > 60) {
    errors.name = "Name must be between 3 and 60 characters in length";
  }

  if (!email || !validator.isEmail(email)) {
    errors.email = "A valid email address is required";
  }

  if (!password || password.length < 6) {
    errors.password = "Password must be at least 6 characters";
  }

  return errors;
};

const asyncValidate = (values, dispatch) => {

  const checkName = () => {
    if (!values.name) return false;
    return api.isName(values.name)
    .then(result => {
      if (result.data) {
        return { name: "This name is already in use" };
      }
  })};

  const checkEmail = () => {
    if (!values.email) return false;
    return api.isEmail(values.email)
    .then(result => {
      if (result.data) {
        return { email: "This email is already in use" };
      }
  })};

  return Promise.all([
    checkEmail(),
    checkName()
  ])
  .then(errors => {
    const result = _.reduce(errors, (result, error) => {
      if (error) {
        return Object.assign({}, result, error);
      }
    }, {});
    return result;
  });

};

export class Signup extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindActionCreators(auth, dispatch);
  }

  handleSubmit(values) {
    const { name, email, password } = values;
    this.actions.signup(name, email, password);
  }

  render() {

    const {
      fields: { name, email, password },
      error,
      handleSubmit,
      resetForm,
      submitting
    } = this.props;

    const onSubmit = handleSubmit(this.handleSubmit.bind(this));

    return (
    <DocumentTitle title={getTitle('Signup')}>
      <div>
        <h2>Join PodBaby today.</h2>
        <hr />
        <p className="lead">
          As a member you can subscribe to podcast feeds and keep track of your favorite episodes.
        </p>
        <form className="form-horizontal"
              onSubmit={onSubmit}>
              <Input hasFeedback={name.touched} bsStyle={name.touched ? ( name.error ? 'error': 'success' ) : undefined}>
                <input type="text" className="form-control" placeholder="Name" {...name} />
                {name.touched && name.error && <div className="help-block">{name.error}</div>}
              </Input>
              <Input hasFeedback={email.touched} bsStyle={email.touched ? ( email.error ? 'error': 'success' ) : undefined}>
                <input type="email" className="form-control" placeholder="Email address" {...email} />
                {email.touched && email.error && <div className="help-block">{email.error}</div>}
              </Input>
              <Input hasFeedback={password.touched} bsStyle={password.touched ? ( password.error ? 'error': 'success' ) : undefined}>
                <input type="password" className="form-control" placeholder="Password" {...password} />
                {password.touched && password.error && <div className="help-block">{password.error}</div>}
              </Input>
            <Button
              bsStyle="primary"
              disabled={submitting}
              onClick={onSubmit}
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
  fields: PropTypes.object.isRequired,
  handleSubmit: PropTypes.func.isRequired,
  error: PropTypes.string,
  resetForm: PropTypes.func.isRequired,
  submitting: PropTypes.bool.isRequired,
  asyncValidating: PropTypes.string.isRequired,
  dispatch: PropTypes.func.isRequired
};

const fields = ['name', 'email', 'password'];
const asyncBlurFields = ['name', 'email'];

export default reduxForm({
  form: 'signup',
  fields,
  validate,
  asyncValidate,
  asyncBlurFields })(Signup);
