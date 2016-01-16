import React, { PropTypes } from 'react';

import { reduxForm } from 'redux-form';


import {
  Button,
  ButtonGroup,
  Modal,
} from 'react-bootstrap';

import * as api from '../api';
import Icon from '../components/icon';
import { FormGroup } from '../components/form';

const validateRecoverPassword = values => {
  return values.identifier ? {} : { identifier: 'You must provide a name or email' };
};


export class RecoverPasswordModal extends React.Component {

  handleSubmit(values) {
    const { identifier } = values;
    const { resetForm, onComplete } = this.props;

    return new Promise((resolve, reject) => {
      return api.recoverPassword(identifier)
      .then(() => {
        onComplete();
        resetForm();
        resolve();
      }, error => {
        reject(error.data);
      });
    });
  }

  render() {
    const {
      fields: { identifier },
      handleSubmit,
      resetForm,
      submitting,
      show,
      onClose,
      container,
    } = this.props;

    const handleClose = () => {
      resetForm();
      onClose();
    };

    return (
      <Modal
        show={show}
        aria-labelledby="recover-password-modal-title"
        container={container}
        onHide={onClose}
      >
        <Modal.Header closeButton>
          <Modal.Title id="recover-password-modal-title">Recover password</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <p>We'll send you a new random password so you can log back in again.</p>
          <form className="form" onSubmit={handleSubmit(this.handleSubmit.bind(this))}>
            <FormGroup field={identifier}>
              <input
                type="text"
                className="form-control"
                placeholder="Email address or name"
                {...identifier}
              />
            </FormGroup>
            <ButtonGroup>
              <Button
                bsStyle="primary"
                disabled={submitting}
                type="submit"
              ><Icon icon="send" /> Send
              </Button>
              <Button bsStyle="default" onClick={handleClose}><Icon icon="remove" /> Cancel</Button>
            </ButtonGroup>
          </form>
        </Modal.Body>
      </Modal>
    );
  }

}

RecoverPasswordModal.propTypes = {
  resetForm: PropTypes.func.isRequired,
  onComplete: PropTypes.func.isRequired,
  handleSubmit: PropTypes.func.isRequired,
  onClose: PropTypes.func.isRequired,
  fields: PropTypes.object.isRequired,
  submitting: PropTypes.bool.isRequired,
  show: PropTypes.bool.isRequired,
  container: PropTypes.object.isRequired,
};

export default reduxForm({
  form: 'recover-password',
  fields: ['identifier'],
  validate: validateRecoverPassword,
})(RecoverPasswordModal);
