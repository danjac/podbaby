import _ from 'lodash';
import React, { PropTypes } from 'react';
import { Link } from 'react-router';

import {
  ButtonGroup,
  Button,
} from 'react-bootstrap';

import { isMobile } from './utils';
import Image from './image';
import Icon from './icon';

class Player extends React.Component {

  constructor(props) {
    super(props);
    this.handleClose = this.handleClose.bind(this);
    this.handleTimeUpdate = this.handleTimeUpdate.bind(this);
    this.handleBookmark = this.handleBookmark.bind(this);
    this.handlePlay = this.handlePlay.bind(this);
    this.handlePlayNext = this.handlePlayNext.bind(this);
    this.handlePlayLast = this.handlePlayLast.bind(this);
    this.handlePlayRandom = this.handlePlayRandom.bind(this);
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

  handlePlayRandom() {
    this.props.onPlayRandom();
  }

  handlePlayNext() {
    this.props.onPlayNext();
  }

  handlePlayLast() {
    this.props.onPlayLast();
  }

  renderButtons(title) {
    const btnStyle = {
      color: '#fff',
      backgroundColor: '#222',
    };

    const { podcast } = this.props.player;

    const showPlaylistButtons = (
      this.props.isLoggedIn &&
      this.props.bookmarks && this.props.bookmarks.length > 0
    );

    return (
        <ButtonGroup style={{ color: '#fff' }}>
          {showPlaylistButtons ?
          <span>
            <Button
              title="Play previous podcast in my bookmarks"
              style={btnStyle}
              onClick={this.handlePlayLast}
            >
              <Icon icon="backward" />
            </Button>
            <Button
              title="Play next podcast in my bookmarks"
              style={btnStyle}
              onClick={this.handlePlayNext}
            >
              <Icon icon="forward" />
            </Button>
            <Button
              title="Play random podcast from my bookmarks"
              style={btnStyle}
              onClick={this.handlePlayRandom}
            >
              <Icon icon="random" />
            </Button>
           </span>
          : ''}
          <Button
            title="Close player"
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
    const maxLength = isMobile() ? 30 : 200;
    const title = _.truncate(fullTitle, { length: maxLength });

    return (
      <div className="container text-center" style={{
        position: 'fixed',
        padding: 5,
        opacity: 0.8,
        backgroundColor: '#222',
        color: '#fff',
        fontWeight: 'bold',
        height: 100,
        bottom: 0,
        width: '100%',
        left: 0,
        right: 0,
        zIndex: 100,
      }}
      >
      <div className="media">
        <div className="media-left media-middle">
          <Link to={`/channel/${podcast.channelId}/`}>
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
         </Link>
        </div>
        <div className="media-body">
          <div>
            <b><Link
              style={{ color: '#fff' }}
              title={fullTitle}
              to={`/podcast/${podcast.id}/`}
            >{title}</Link></b><br />
            {this.renderButtons(fullTitle)}
          </div>
          </div>
         </div>
          <div>
            <audio
              controls
              autoPlay
              onPlay={this.handlePlay}
              onTimeUpdate={this.handleTimeUpdate}
              src={podcast.enclosureUrl}
              style={{
                backgroundColor: '#222',
                color: '#fff',
                width: '100%',
              }}
            >
              <source src={podcast.enclosureUrl} />
              Download from <a download href={podcast.enclosureUrl}>here</a>.
            </audio>
            </div>
        </div>
    );
  }
}


Player.propTypes = {
  onClose: PropTypes.func.isRequired,
  onTimeUpdate: PropTypes.func.isRequired,
  onToggleBookmark: PropTypes.func.isRequired,
  onPlayNext: PropTypes.func.isRequired,
  onPlayLast: PropTypes.func.isRequired,
  onPlayRandom: PropTypes.func.isRequired,
  player: PropTypes.object.isRequired,
  isLoggedIn: PropTypes.bool.isRequired,
  bookmarks: PropTypes.array,
};

export default Player;
