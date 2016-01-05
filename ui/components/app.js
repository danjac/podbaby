import React, { PropTypes } from 'react';
import { Link } from 'react-router';

import 'bootswatch/spacelab/bootstrap.min.css';
import 'font-awesome/css/font-awesome.min.css';

import {
  Badge,
  Alert,
  Grid,
  Row,
  Col
} from 'react-bootstrap';

import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { pushPath } from 'redux-simple-router';

import * as actions from '../actions';

import Player from './player';
import NavBar from './navbar';
import AddChannelModal from './add_channel';


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
      bookmarks: bindActionCreators(actions.bookmarks, dispatch),
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
    this.actions.player.close(this.props.player);
  }

  handleTogglePlayerBookmark() {
    if (this.props.player.podcast) {
      this.actions.bookmarks.addBookmark(this.props.player.podcast.id);
    }
  }

  handleUpdatePlayerTime(event) {
    this.actions.player.updateTime(
      this.props.player,
      event.currentTarget.currentTime,
    );
  }

  render() {

    const { createHref } = this.props.history;
    const { isLoggedIn } = this.props.auth;

    const pageContent = (
        <div className="container">
          {this.props.children}
        </div>
    );

    const alertList = (
        <AlertList alerts={this.props.alerts}
                   onDismissAlert={this.handleDismissAlert.bind(this)} />
    );
    if (isLoggedIn) {
      return (
        <div>
          <NavBar onLogout={this.handleLogout.bind(this)}
                  onOpenAddChannelForm={this.handleOpenAddChannelForm.bind(this)}
                   {...this.props} />
          {pageContent}
          {this.props.player.isPlaying ?
            <Player player={this.props.player}
                    onToggleBookmark={this.handleTogglePlayerBookmark.bind(this)}
                    onTimeUpdate={this.handleUpdatePlayerTime.bind(this)}
                    onClosePlayer={this.handleClosePlayer.bind(this)}/> : ''}

          {alertList}
          <AddChannelModal {...this.props.addChannel}
                           container={this}
                           onAdd={this.handleAddChannel.bind(this)}
                           onClose={this.handleCloseAddChannelForm.bind(this)} />
        </div>
      );
    } else {
      return (
        <div>
          {pageContent}
          {alertList}
        </div>
      );
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
