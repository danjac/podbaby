import React, { PropTypes } from 'react';
import { connect } from 'react-redux';

import 'bootswatch/paper/bootstrap.min.css';
import 'font-awesome/css/font-awesome.min.css';

import { Alert } from 'react-bootstrap';

import * as actions from '../actions';
import { bindAllActionCreators } from '../actions/utils';

import Player from '../components/player';
import NavBar from '../components/navbar';
import AddChannelModal from '../components/add_channel';


const AlertList = props => {
  if (props.alerts.length === 0) return <div></div>;

  return (
    <div
      className="container"
      style={{
        position: 'fixed',
        height: '50px',
        width: '500px',
        opacity: 0.9,
        textAlign: 'center',
        margin: '5% auto',
        left: 0,
        right: 0,
        zIndex: 200,
      }}
    >
    {props.alerts.map(alert => {
      const dismissAlert = () => props.onDismissAlert(alert.id);
      return (
      <Alert key={alert.id} bsStyle={alert.status} onDismiss={dismissAlert} dismissAfter={3000}>
        <p><b style={{ opacity: 1.0 }}>{alert.message}</b></p>
      </Alert>);
    })}
    </div>
  );
};

AlertList.propTypes = {
  alerts: PropTypes.array.isRequired,
};


export class App extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindAllActionCreators(actions, dispatch);

    this.handleLogout = this.handleLogout.bind(this);
    this.handleDismissAlert = this.handleDismissAlert.bind(this);
    this.handleOpenAddChannelForm = this.handleOpenAddChannelForm.bind(this);
    this.handleCloseAddChannelForm = this.handleCloseAddChannelForm.bind(this);
    this.handleAddChannelComplete = this.handleAddChannelComplete.bind(this);
    this.handleClosePlayer = this.handleClosePlayer.bind(this);
    this.handleTogglePlayerBookmark = this.handleTogglePlayerBookmark.bind(this);
    this.handleUpdatePlayerTime = this.handleUpdatePlayerTime.bind(this);
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

  handleAddChannelComplete(channel) {
    this.actions.addChannel.complete(channel);
  }

  handleDismissAlert(id) {
    this.actions.alerts.dismissAlert(id);
  }

  handleClosePlayer() {
    this.actions.player.close();
  }

  handleTogglePlayerBookmark() {
    if (this.props.player.podcast) {
      this.actions.bookmarks.toggleBookmark(this.props.player.podcast);
    }
  }

  handleUpdatePlayerTime(event) {
    this.actions.player.updateTime(
      event.currentTarget.currentTime
    );
  }

  render() {
    const { isActive } = this.props.history;
    const { isLoggedIn } = this.props.auth;

    const hideNavbar = isActive('/front/');

    const pageContent = (
        <div className="container">
          {this.props.children}
        </div>
    );

    const alertList = (
      <AlertList
        alerts={this.props.alerts}
        onDismissAlert={this.handleDismissAlert}
      />
    );

    return (
      <div>
        {hideNavbar ? '' :
        <NavBar
          onLogout={this.handleLogout}
          onOpenAddChannelForm={this.handleOpenAddChannelForm}
          {...this.props}
        />}
        {alertList}
        {pageContent}
        {this.props.player.isPlaying ?
        <Player
          player={this.props.player}
          isLoggedIn={isLoggedIn}
          onToggleBookmark={this.handleTogglePlayerBookmark}
          onTimeUpdate={this.handleUpdatePlayerTime}
          onClose={this.handleClosePlayer}
        /> : ''}

        <AddChannelModal
          {...this.props.addChannel}
          container={this}
          onComplete={this.handleAddChannelComplete}
          onClose={this.handleCloseAddChannelForm}
        />
    </div>
    );
  }
}


App.propTypes = {
  dispatch: PropTypes.func.isRequired,
  routing: PropTypes.object.isRequired,
  history: PropTypes.object.isRequired,
  children: PropTypes.array.isRequired,
  auth: PropTypes.object.isRequired,
  addChannel: PropTypes.object.isRequired,
  player: PropTypes.object.isRequired,
  alerts: PropTypes.array.isRequired,
};


const mapStateToProps = state => {
  const { routing, auth, addChannel, player, alerts } = state;
  return {
    routing,
    auth,
    addChannel,
    player,
    alerts,
  };
};

export default connect(mapStateToProps)(App);
