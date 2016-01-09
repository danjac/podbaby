import React from 'react';
import { reduxForm } from 'redux-form';
import validator from 'validator';

import {
  Modal,
  Input,
  Button,
  ButtonGroup,
  ProgressBar
} from 'react-bootstrap';

import Icon from './icon';
import { FormGroup } from './form';

const validate = values => {
  return values.url && validator.isURL(values.url) ? {} : {
    url: "You must provide a valid URL"
  };
}

export class AddChannelModal extends React.Component {

  constructor(props) {
    super(props);
    this.state = this.getDefaultState();
  }

  getDefaultState() {
    return {
      progress: 0,
      interval: null
    }
  }

  componentWillReceiveProps(newProps) {
    if (newProps.pending && !this.props.pending) {
      this.setState({
        interval: window.setInterval(() => {
          this.setState({ progress: this.state.progress + 1 });
        }, 100)
      });
    } else if (!newProps.pending && this.props.pending) {
      window.clearInterval(this.state.interval);
      this.setState(this.getDefaultState());
    }
    return this.props !== newProps;
  }

  render() {
    const { show, onClose, container, pending } = this.props;

    const {
      handleSubmit,
      fields: { url },
      submitting,
      resetForm
    } = this.props;

    const handleAdd = values => {
      this.props.onAdd(values.url);
      resetForm();
    };

    const handleClose = event => {
      resetForm();
      onClose(event);
    };

    return (
      <Modal show={show}
             aria-labelledby="add-channel-modal-title"
             container={container}
             onHide={handleClose}>
        <Modal.Header closeButton>
          <Modal.Title id="add-channel-modal-title">Add a new channel</Modal.Title>
        </Modal.Header>
        <Modal.Body>
            {pending ? (
            <div>
              <ProgressBar now={this.state.progress} />
            </div>
            ) : (
            <form className="form" onSubmit={handleSubmit(handleAdd)}>
               <FormGroup field={url}>
                <input type="text" className="form-control"  {...url} />
              </FormGroup>
            <p>Enter the URL of the RSS feed you want to subscribe to, for example:
              <br /><em>http://joeroganexp.joerogan.libsynpro.com/rss</em>
            </p>
              <ButtonGroup>
              <Button bsStyle="primary" type="submit"><Icon icon="plus" /> Add channel</Button>
              <Button bsStyle="default" onClick={handleClose}><Icon icon="remove" /> Cancel</Button>
            </ButtonGroup>
            </form>
            )}
        </Modal.Body>
      </Modal>
    );
  }

}


export default reduxForm({
  form: 'add-channel',
  fields: ['url'],
  validate
})(AddChannelModal);
