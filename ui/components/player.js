import _ from 'lodash';
import React, { PropTypes } from 'react';
import { Link } from 'react-router';

import {
  ButtonGroup,
  Button,
} from 'react-bootstrap';

import Image from './image';
import Icon from './icon';

class Player extends React.Component {

  constructor(props) {
    super(props);
    this.handleClose = this.handleClose.bind(this);
    this.handleTimeUpdate = this.handleTimeUpdate.bind(this);
    this.handleBookmark = this.handleBookmark.bind(this);
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

  renderButtons(title) {
    const btnStyle = {
      color: '#fff',
      backgroundColor: '#222',
    };

    const { podcast } = this.props.player;

    return (
        <ButtonGroup style={{ color: '#fff', float: 'right' }}>
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
            title={`Download ${title}`}
            className="btn btn-default"
            style={btnStyle}
            href={podcast.enclosureUrl}
          ><Icon icon="download" /></a>
           {this.props.isLoggedIn ?
           <Button
             title={podcast.isBookmarked ? 'Remove bookmark' : 'Add bookmark '}
             pullRight
             style={btnStyle}
             onClick={this.handleBookmark}
           ><Icon icon={podcast.isBookmarked ? 'bookmark' : 'bookmark-o'} />
          </Button> : ''}
        </ButtonGroup>
    );
  }

  render() {
    const { player } = this.props;
    const { podcast } = player;
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
        height: 60,
        bottom: 0,
        width: '100%',
        left: 0,
        right: 0,
        zIndex: 100,
      }}
      >
      <div className="media">
        <div className="media-left media-middle">
          <Image
            className="media-object"
            src={podcast.image}
            errSrc="/static/podcast.png"
            imgProps={{
              height: 40,
              width: 40,
              alt: podcast.name,
            }}
          />
        </div>
        <div className="media-body">
          <div>
            <b><Link
              style={{ color: '#fff' }}
              title={fullTitle}
              to={`/podcast/${podcast.id}/`}
            >{title}</Link></b>
            {this.renderButtons(fullTitle)}
          </div>
          <div>
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
            </div>
          </div>
         </div>
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
