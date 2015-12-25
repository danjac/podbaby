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

class MainNav extends React.Component {

  search(event) {
    event.preventDefault();
    const node = this.refs.search.getInputDOMNode();
    this.props.search(node.value);
    node.value = "";
  }

  render() {

    const { auth } = this.props;
    const { createHref, isActive } = this.props.history;
    const searchIcon = <Glyphicon glyph="search" />;

    return (
      <Navbar fixedTop={true}>
        <Navbar.Header>
          <Navbar.Brand>
            <Link to={auth.isLoggedIn ? "/podcasts/new/" : "/" }><Glyphicon glyph="headphones" /> PodBaby</Link>
          </Navbar.Brand>
        </Navbar.Header>

        {auth.isLoggedIn ?
        <form className="navbar-form navbar-left" role="search" onSubmit={this.search.bind(this)}>
          <Input ref="search" type="search" placeholder="Find podcast or channel" addonBefore={searchIcon} />
        </form> : ''}

        {auth.isLoggedIn ?

        <Nav pullLeft={true}>
          <NavItem active={isActive("/podcasts/new/")}
            href={createHref("/podcasts/new/")}><Glyphicon glyph="flash" /> New podcasts <Badge>24</Badge></NavItem>
          <NavItem active={isActive("/podcasts/subscriptions/")}
                   href={createHref("/podcasts/subscriptions/")}><Glyphicon glyph="list" /> Subscriptions</NavItem>
          <NavItem href="#"><Glyphicon glyph="pushpin" /> Bookmarks</NavItem>
          <NavItem onClick={this.props.openAddChannelForm} href="#"><Glyphicon glyph="plus" /> Add new</NavItem>
        </Nav> : ''}

        {auth.isLoggedIn ?
        <Nav pullRight={true}>
          <NavItem href="#"><Glyphicon glyph="cog" /> {auth.name}</NavItem>
          <NavItem href="#" onClick={this.props.logout}><Glyphicon glyph="log-out" /> Logout</NavItem>
        </Nav> :
        <Nav pullRight={true}>
          <NavItem active={isActive("/login/")}
                   href={createHref("/login/")}><Glyphicon glyph="log-in" /> Login</NavItem>
          <NavItem active={isActive("/signup/")}
                   href={createHref("/signup/")}><Glyphicon glyph="user" /> Signup</NavItem>
        </Nav>}

      </Navbar>
    );
  }
}


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
            <form className="form" onSubmit={close}>
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

export class Player extends React.Component {
  render() {
    const { podcast } = this.props.player;
    return (
    <p>
      Currently playing: <b>{podcast.name}</b><br />
      <audio controls={true} autoPlay={true} src={podcast.enclosureUrl}>
        <source src={podcast.enclosureUrl} />
        Download from <a href="#">here</a>.
      </audio>
    </p>
    );
  }
}


export class App extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;

    this.actions = {
      auth: bindActionCreators(actions.auth, dispatch),
      addChannel: bindActionCreators(actions.addChannel, dispatch),
      search: bindActionCreators(actions.search, dispatch)
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

  search(query) {
    this.actions.search.search(query);
  }

  render() {

    return (
      <div>
        <MainNav logout={this.logout.bind(this)}
                 openAddChannelForm={this.openAddChannelForm.bind(this)}
                 search={this.search.bind(this)}
                 {...this.props} />
        <div className="container" style={ { marginTop: "80px", opacity: 0.5, backgroundColor: "#eee" } }>
          {this.props.auth.isLoggedIn && this.props.player.isPlaying ? <Player player={this.props.player} /> : ''}
        </div>
        <div className="container" style={ { marginTop: "80px" } }>
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
  dispatch: PropTypes.func.isRequired,
  routing: PropTypes.object.isRequired,
  auth: PropTypes.object
};


const mapStateToProps = state => {
  const { routing, auth, addChannel, player } = state;
  return {
    routing,
    auth,
    addChannel,
    player
  };
};

export default connect(mapStateToProps)(App);
