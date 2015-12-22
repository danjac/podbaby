import React, { PropTypes } from 'react';
import { Link } from 'react-router';

import 'bootswatch/paper/bootstrap.min.css';

import {
  Nav,
  NavItem,
  Navbar,
  Glyphicon,
  Badge,
  Input,
  Button,
  ButtonGroup,
  Modal
} from 'react-bootstrap';

import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

import * as actions from '../actions';

const MainNav = props => {
  const { auth } = props;
  const { createHref, isActive } = props.history;
  const searchIcon = <Glyphicon glyph="search" />;
  return (
    <Navbar fixedTop={true}>
      <Navbar.Header>
        <Navbar.Brand>
          <Link to={auth.isLoggedIn ? "/podcasts/" : "/" }><Glyphicon glyph="headphones" /> PodBaby</Link>
        </Navbar.Brand>
      </Navbar.Header>

      {auth.isLoggedIn ?

      <Nav pullLeft={true}>
        <NavItem active={isActive("/podcasts/")} href={createHref("/podcasts/")}><Glyphicon glyph="flash" /> New podcasts <Badge>24</Badge></NavItem>
        <NavItem active={isActive("/podcasts/list/")} href="#"><Glyphicon glyph="list" /> Subscriptions</NavItem>
        <NavItem href="#"><Glyphicon glyph="pushpin" /> Pins <Badge>4</Badge></NavItem>
        <NavItem onClick={props.openAddChannelForm} href="#"><Glyphicon glyph="plus" /> Add new podcast</NavItem>
      </Nav> : ''}

      {auth.isLoggedIn ?
      <form className="navbar-form navbar-left" role="search">
        <Input type="search" placeholder="Find podcast or channel" addonBefore={searchIcon}/>
      </form> : ''}

      {auth.isLoggedIn ?
      <Nav pullRight={true}>
        <NavItem href="#"><Glyphicon glyph="cog" /> Settings</NavItem>
        <NavItem href="#" onClick={props.logout}><Glyphicon glyph="log-out" /> Logout</NavItem>
      </Nav> :
      <Nav pullRight={true}>
        <NavItem active={isActive("/login/")}
                 href={createHref("/login/")}><Glyphicon glyph="log-in" /> Login</NavItem>
        <NavItem active={isActive("/signup/")}
                 href={createHref("/signup/")}><Glyphicon glyph="user" /> Signup</NavItem>
      </Nav>}

    </Navbar>
  );
};


class AddChannelModal extends React.Component {

  render() {
    const { show, close, container } = this.props;
    return (
      <Modal show={show}
             aria-labelledby="add-channel-modal-title"
             container={container}
             onHide={close}>
        <Modal.Header closeButton>
          <Modal.Title id="add-channel-modal-title">Add a new channel</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <form className="form form-horizontal" onSubmit={close}>
            <Input required type="text" placeholder="RSS URL of the channel" />
            <ButtonGroup>
            <Button bsStyle="primary" type="submit"><Glyphicon glyph="plus" /> Add channel</Button>
            <Button bsStyle="default" onClick={close}><Glyphicon glyph="remove" /> Cancel</Button>
          </ButtonGroup>
          </form>
        </Modal.Body>
      </Modal>
    );
  }


}

export class App extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;

    this.actions = {
      auth: bindActionCreators(actions.auth, dispatch),
      addChannel: bindActionCreators(actions.addChannel, dispatch)
    }
  }

  logout(event) {
    event.preventDefault();
    this.actions.auth.logout();
  }

  openAddChannelForm(event) {
    event.preventDefault();
    this.actions.addChannel.open();
  }

  closeAddChannelForm(event) {
    event.preventDefault();
    this.actions.addChannel.close();
  }

  render() {

    return (
      <div>
        <MainNav logout={this.logout.bind(this)}
                 openAddChannelForm={this.openAddChannelForm.bind(this)}
                 {...this.props} />
        <div className="container" style={ { marginTop: "80px"}  }>
          {this.props.children}
        </div>
        <AddChannelModal show={this.props.addChannel.show}
                         container={this}
                         close={this.closeAddChannelForm.bind(this)} />
      </div>
    );
  }
}

App.propTypes = {
  routing: PropTypes.object.isRequired,
  dispatch: PropTypes.func.isRequired,
  auth: PropTypes.object
};


const mapStateToProps = state => {
  const { routing, auth, addChannel } = state;
  return {
    routing,
    auth,
    addChannel
  };
};

export default connect(mapStateToProps)(App);
