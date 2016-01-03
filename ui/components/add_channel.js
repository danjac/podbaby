import React from 'react';

import {
  Modal,
  Input,
  Button,
  Glyphicon,
  ButtonGroup,
  ProgressBar
} from 'react-bootstrap';

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

  handleAdd(event){
    event.preventDefault();
    const node = this.refs.url.getInputDOMNode();
    this.props.onAdd(node.value);
    node.value = "";
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
    return (
      <Modal show={show}
             aria-labelledby="add-channel-modal-title"
             container={container}
             onHide={onClose}>
        <Modal.Header closeButton>
          <Modal.Title id="add-channel-modal-title">Add a new channel</Modal.Title>
        </Modal.Header>
        <Modal.Body>
            {pending ? (
            <div>
              <ProgressBar now={this.state.progress} />
            </div>
            ) : (
            <form className="form" onSubmit={this.handleAdd.bind(this)}>
              <Input required type="text" placeholder="RSS URL of the channel" ref="url" />
              <ButtonGroup>
              <Button bsStyle="primary" type="submit"><Glyphicon glyph="plus" /> Add channel</Button>
              <Button bsStyle="default" onClick={onClose}><Glyphicon glyph="remove" /> Cancel</Button>
            </ButtonGroup>
            </form>
            )}
        </Modal.Body>
      </Modal>
    );
  }

}

export default AddChannelModal;
