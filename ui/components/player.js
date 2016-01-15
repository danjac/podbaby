import _ from 'lodash';
import React, { PropTypes } from 'react';
import { Link } from 'react-router';

import {
  Grid,
  Row,
  Col,
  ButtonGroup,
  Button,
} from 'react-bootstrap';

import Icon from './icon';

class Player extends React.Component {

  constructor(props) {
    super(props);
    this.handleClose = this.handleClose.bind(this);
    this.handleTimeUpdate = this.handleTimeUpdate.bind(this);
    this.handlePlay = this.handlePlay.bind(this);
  }

  handleClose(event) {
    event.preventDefault();
    this.props.onClose();
  }

  handleTimeUpdate(event) {
    this.props.onTimeUpdate(event);
  }

  handlePlay(event) {
    const { currentTarget } = event;
    currentTarget.currentTime = this.props.player.currentTime;
  }

  handleBookmark() {
    this.props.onToggleBookmark();
  }

  render() {
    const { player, isLoggedIn } = this.props;
    const { podcast } = player;
    const btnStyle = {
      color: '#fff',
      backgroundColor: '#222',
    };

    const fullTitle = podcast.name + ' : ' + podcast.title;
    const title = _.truncate(fullTitle, 50);

    return (
      <div className="container" style={{
        position: 'fixed',
        padding: 5,
        opacity: 0.8,
        backgroundColor: '#222',
        color: '#fff',
        fontWeight: 'bold',
        height: 40,
        bottom: 0,
        width: '100%',
        left: 0,
        right: 0,
        zIndex: 100,
      }}
      >
        <Grid>
          <Row>
            <Col xs={6} sm={6} md={4}>
              <audio
                controls
                autoPlay
                onPlay={this.handlePlay}
                onTimeUpdate={this.handleTimeUpdate}
                src={podcast.enclosureUrl}
              >
                <source src={podcast.enclosureUrl} />
                Download from <a download href={podcast.enclosureUrl}>here</a>.
              </audio>
            </Col>
            <Col md={4} className="hidden-xs hidden-sm">
              <b><Link
                style={{ color: '#fff' }}
                title={fullTitle}
                to={`/podcast/${podcast.id}/`}
              >{title}</Link></b>
            </Col>
            <Col xs={6} sm={6} md={4} mdPush={2} xsPush={1}>
              <ButtonGroup style={{ color: '#fff' }}>
                <Button
                  title="Close player"
                  pullRight
                  style={btnStyle}
                  onClick={this.handleClose}
                >
                  <Icon icon="stop" />
                </Button>
                <a
                  download
                  title={`Download ${fullTitle}`}
                  className="btn btn-default"
                  style={btnStyle}
                  href={podcast.enclosureUrl}
                ><Icon icon="download" /></a>
                 {isLoggedIn ?
                 <Button
                   title={podcast.isBookmarked ? 'Remove bookmark' : 'Add bookmark '}
                   pullRight
                   style={btnStyle}
                   onClick={this.handleBookmark}
                 ><Icon icon={podcast.isBookmarked ? 'bookmark' : 'bookmark-o'} />
                </Button> : ''}
              </ButtonGroup>
            </Col>
          </Row>
        </Grid>
    </div>
    );
  }
}


Player.propTypes = {
  onClose: PropTypes.func.isRequired,
  onTimeUpdate: PropTypes.func.isRequired,
  onToggleBookmark: PropTypes.func.isRequired,
  player: PropTypes.object.isRequired,
  isLoggedIn: PropTypes.bool.isRequired,
};

export default Player;
