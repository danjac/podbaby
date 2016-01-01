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
  Modal,
  Alert,
  Grid,
  Row,
  Col
} from 'react-bootstrap';

import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { pushPath } from 'redux-simple-router';

import * as actions from '../actions';

class MainNav extends React.Component {

  render() {

    const { auth } = this.props;
    const { createHref, isActive } = this.props.history;

    return (
      <Navbar fixedTop={true}>
        <Navbar.Header>
          <Navbar.Brand>
            <Link style={{ fontFamily: "GoodDog" }} to="/podcasts/new/"><Glyphicon glyph="headphones" /> PodBaby</Link>
          </Navbar.Brand>
        </Navbar.Header>


        <Nav pullLeft={true}>
          <NavItem active={isActive("/podcasts/search/")}
                   href={createHref("/podcasts/search/")}><Glyphicon glyph="search" /> Discover</NavItem>
          <NavItem active={isActive("/podcasts/new/")}
                   href={createHref("/podcasts/new/")}><Glyphicon glyph="flash" /> New episodes</NavItem>
          <NavItem active={isActive("/podcasts/subscriptions/")}
                   href={createHref("/podcasts/subscriptions/")}><Glyphicon glyph="list" /> Subscriptions</NavItem>
          <NavItem active={isActive("/podcasts/bookmarks/")}
                   href={createHref("/podcasts/bookmarks/")}><Glyphicon glyph="bookmark" /> Bookmarks</NavItem>
          <NavItem onClick={this.props.onOpenAddChannelForm} href="#"><Glyphicon glyph="plus" /> Add new</NavItem>
        </Nav>

        <Nav pullRight={true}>
          <NavItem active={isActive("/user/")}
                    href={createHref("/user/")}><Glyphicon glyph="cog" /> {auth.name}</NavItem>
          <NavItem href="#" onClick={this.props.onLogout}><Glyphicon glyph="log-out" /> Logout</NavItem>
        </Nav>

      </Navbar>
    );
  }
}


class AddChannelModal extends React.Component {

  handleAdd(event){
    event.preventDefault();
    const node = this.refs.url.getInputDOMNode();
    this.props.onAdd(node.value);
    node.value = "";
  }

  render() {
    const { show, onClose, container } = this.props;
    return (
      <Modal show={show}
             aria-labelledby="add-channel-modal-title"
             container={container}
             onHide={onClose}>
        <Modal.Header closeButton>
          <Modal.Title id="add-channel-modal-title">Add a new channel</Modal.Title>
        </Modal.Header>
        <Modal.Body>
            <form className="form" onSubmit={this.handleAdd.bind(this)}>
              <Input required type="text" placeholder="RSS URL of the channel" ref="url" />
              <ButtonGroup>
              <Button bsStyle="primary" type="submit"><Glyphicon glyph="plus" /> Add channel</Button>
              <Button bsStyle="default" onClick={onClose}><Glyphicon glyph="remove" /> Cancel</Button>
            </ButtonGroup>
            </form>
        </Modal.Body>
      </Modal>
    );
  }


}

export class Player extends React.Component {
  handleClose(event) {
    event.preventDefault();
    this.props.onClosePlayer();
  }

  render() {
    const { podcast } = this.props.player;
    return (
      <footer style={{
        position:"fixed",
        padding: "5px",
        opacity: 0.7,
        backgroundColor: "#222",
        color: "#fff",
        fontWeight: "bold",
        height: "50px",
        bottom: 0,
        width: "100%",
        zIndex: 100
        }}>
        <Grid>
          <Row>
            <Col xs={6} md={6}>
              Currently playing: <b><Link to={`/podcasts/channel/${podcast.channelId}/`}>{podcast.name}</Link> : {podcast.title}</b>
            </Col>
            <Col xs={3} md={4}>
              <audio controls={true} autoPlay={true} src={podcast.enclosureUrl}>
                <source src={podcast.enclosureUrl} />
                Download from <a href="#">here</a>.
              </audio>
            </Col>
            <Col xs={3} md={2} mdPush={2}>
              <a href="#" onClick={this.handleClose.bind(this)} style={{ color: "#fff" }}><Glyphicon glyph="remove" /> Close</a>
            </Col>
          </Row>
        </Grid>
    </footer>
    );
  }
}


const AlertList = props => {

  if (props.alerts.length === 0) return <div></div>;

  return (
    <div className="container" style={{
        position: "fixed",
        height: "50px",
        width: "100%",
        bottom: 20,
        zIndex: 200
      }}>
      {props.alerts.map(alert => {
        const dismissAlert = () => props.onDismissAlert(alert.id);
        return (<Alert key={alert.id} bsStyle={alert.status} onDismiss={dismissAlert} dismissAfter={3000}>
          <p>{alert.message}</p>
        </Alert>);
      })}
    </div>
  );
};

export class App extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;

    this.actions = {
      auth: bindActionCreators(actions.auth, dispatch),
      addChannel: bindActionCreators(actions.addChannel, dispatch),
      search: bindActionCreators(actions.search, dispatch),
      player: bindActionCreators(actions.player, dispatch),
      alerts: bindActionCreators(actions.alerts, dispatch)
    }
  }

  handleLogout(event) {
    event.preventDefault();
    this.actions.auth.logout();
  }

  handleOpenAddChannelForm(event) {
    event.preventDefault();
    this.actions.addChannel.open();
  }

  handleCloseAddChannelForm(event) {
    event.preventDefault();
    this.actions.addChannel.close();
  }

  handleAddChannel(url) {
    this.actions.addChannel.add(url);
  }

  handleDismissAlert(id) {
    this.actions.alerts.dismissAlert(id);
  }

  handleClosePlayer() {
    this.actions.player.setPodcast(null);
  }

  render() {

    const { createHref } = this.props.history;
    const { isLoggedIn } = this.props.auth;

    const pageContent = (
        <div className="container">
          {this.props.children}
        </div>
    );

    if (isLoggedIn) {
      return (
        <div>
          <MainNav onLogout={this.handleLogout.bind(this)}
                   onOpenAddChannelForm={this.handleOpenAddChannelForm.bind(this)}
                   {...this.props} />
          <AlertList alerts={this.props.alerts}
                     onDismissAlert={this.handleDismissAlert.bind(this)} />
          {pageContent}
          {this.props.player.isPlaying ?
            <Player player={this.props.player}
                    createHref={createHref}
                    onClosePlayer={this.handleClosePlayer.bind(this)}/> : ''}

          <AddChannelModal show={this.props.addChannel.show}
                           container={this}
                           onAdd={this.handleAddChannel.bind(this)}
                           onClose={this.handleCloseAddChannelForm.bind(this)} />
        </div>
      );
    } else {
      return pageContent;
    }
  }
}


App.propTypes = {
  dispatch: PropTypes.func.isRequired,
  routing: PropTypes.object.isRequired,
  auth: PropTypes.object,
  addChannel: PropTypes.object,
  player: PropTypes.object,
  alerts: PropTypes.array
};


const mapStateToProps = state => {
  const { routing, auth, addChannel, player, alerts } = state;
  return {
    routing,
    auth,
    addChannel,
    player,
    alerts
  };
};

export default connect(mapStateToProps)(App);
