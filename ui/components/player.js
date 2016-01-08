import React from 'react';
import { Link } from 'react-router';

import {
  Grid,
  Row,
  Col,
  ButtonGroup,
  Button
} from 'react-bootstrap';

import Icon from './icon';

class Player extends React.Component {

  handleClose(event) {
    event.preventDefault();
    this.props.onClosePlayer();
  }

  handleTimeUpdate(event) {
    this.props.onTimeUpdate(event);
  }

  handlePlay(event) {
    event.currentTarget.currentTime = this.props.player.currentTime;
  }

  handleBookmark() {
    this.props.onToggleBookmark();
  }

  render() {
    const { player, isLoggedIn } = this.props;
    const { podcast } = player;
    return (
      <footer style={{
        position:"fixed",
        padding: "5px",
        opacity: 0.8,
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
            <Col xs={6} md={5}>
              <b><Link style={{ color: '#fff' }}
                       to={`/channel/${podcast.channelId}/`}>{podcast.name}</Link> :
                 <Link style={{ color: '#fff' }}
                       to={`/podcast/${podcast.id}/`}>{podcast.title}</Link></b>
            </Col>
            <Col xs={3} md={4}>
              <audio controls
                     autoPlay
                     onPlay={this.handlePlay.bind(this)}
                     onTimeUpdate={this.handleTimeUpdate.bind(this)}
                     src={podcast.enclosureUrl}>
                <source src={podcast.enclosureUrl} />
                Download from <a href="#">here</a>.
              </audio>
            </Col>
            <Col xs={3} md={3} mdPush={2}>
              <ButtonGroup style={{ color: "#fff" }}>
                <a download
                   title="Download this podcast"
                   className="btn btn-sm btn-default"
                   href={podcast.enclosureUrl}><Icon icon="download" /></a>
                 {isLoggedIn ?
                 <Button title={podcast.isBookmarked ? 'Remove bookmark' : 'Add bookmark '}
                         bsSize="sm" onClick={this.handleBookmark.bind(this)}>
                    <Icon icon={podcast.isBookmarked ? 'bookmark' : 'bookmark-o'} />
                </Button> : ''}
                <Button title='Close player' bsSize="sm" onClick={this.handleClose.bind(this)}>
                  <Icon icon="remove" />
                </Button>
              </ButtonGroup>
            </Col>
          </Row>
        </Grid>
    </footer>
    );
  }
}


export default Player;
