import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import DocumentTitle from 'react-document-title';
import { connect } from 'react-redux';

import { Button } from 'react-bootstrap';

import * as actions from '../actions';
import Icon from '../components/icon';
import ChangeEmailForm from '../components/change_email_form';
import ChangePasswordForm from '../components/change_password_form';
import { getTitle } from './utils';


export const DeleteAccountForm = props => {
  return (
      <div>
        <Button
          bsStyle="danger"
          className="form-control"
          onClick={props.onDelete}
        ><Icon icon="trash" /> Delete my account</Button>
        <p className="text-center">
          <b>This will completely and irreversibly delete your account,
            including all your subscriptions and bookmarks.</b>
        </p>
      </div>
  );
};

DeleteAccountForm.propTypes = {
  onDelete: PropTypes.func.isRequired,
};

export class User extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindActionCreators(actions.users, dispatch);
    this.handleDelete = this.handleDelete.bind(this);
    this.handleChangeEmailComplete = this.handleChangeEmailComplete.bind(this);
    this.handleChangePasswordComplete = this.handleChangePasswordComplete.bind(this);
  }

  handleChangeEmailComplete(email) {
    this.actions.changeEmailComplete(email);
  }

  handleChangePasswordComplete() {
    this.actions.changePasswordComplete();
  }

  handleDelete(event) {
    event.preventDefault();
    if (window.confirm(
      'Are you sure you want to delete this account? ' +
      'You will lose all your subscriptions and bookmarks!!!')) {
      this.actions.deleteAccount();
    }
  }

  render() {
    return (
    <DocumentTitle title={getTitle('My settings')}>
      <div>
        <h3>Change my email address</h3>
        <ChangeEmailForm onComplete={this.handleChangeEmailComplete} />
        <h3>Change my password</h3>
        <ChangePasswordForm onComplete={this.handleChangePasswordComplete} />
        <hr />
        <DeleteAccountForm onDelete={this.handleDelete} />
      </div>
    </DocumentTitle>
    );
  }

}

User.propTypes = {
  auth: PropTypes.object.isRequired,
  dispatch: PropTypes.func.isRequired,
};

const mapStateToProps = state => {
  return {
    auth: state.auth,
  };
};

export default connect(mapStateToProps)(User);
