import React, { PropTypes } from 'react';
import { Link } from 'react-router';
import { sanitize, formatPubDate, formatListenDate } from './utils';

import {
  ButtonGroup,
  Button,
  Panel,
  Badge,
} from 'react-bootstrap';


import Icon from './icon';
import Image from './image';


const Buttons = props => {
  const { podcast } = props;
  return (
    <ButtonGroup vertical={props.vertical} style={{ float: 'right' }}>
     <Button
       title={ podcast.isPlaying ? 'Stop' : 'Play' }
       onClick={props.togglePlayer}
     ><Icon icon={ podcast.isPlaying ? 'stop' : 'play' } />
     </Button>
     <Button
       download
       title="Download this podcast"
       className="btn btn-default"
       href={podcast.enclosureUrl}
     ><Icon icon="download" /></Button>
    {props.isLoggedIn ?
    <Button
      onClick={props.toggleBookmark}
      title={podcast.isBookmarked ? 'Remove bookmark' : 'Add to bookmarks'}
    ><Icon icon={podcast.isBookmarked ? 'bookmark' : 'bookmark-o'} />
    </Button> : ''}
    </ButtonGroup>
  );
};

Buttons.propTypes = {
  vertical: PropTypes.bool,
  isLoggedIn: PropTypes.bool.isRequired,
  podcast: PropTypes.object.isRequired,
  toggleBookmark: PropTypes.func.isRequired,
  togglePlayer: PropTypes.func.isRequired,
};

export default function PodcastItem(props) {
  const {
    podcast,
    showChannel,
    showExpanded,
    toggleDetail,
    isLoggedIn,
  } = props;

  const channelUrl = `/channel/${podcast.channelId}/`;
  const podcastUrl = `/podcast/${podcast.id}/`;
  const image = podcast.image || '/static/podcast.png';

  const playedAt = isLoggedIn && podcast.lastPlayedAt ?
    <Badge>Listened {formatListenDate(podcast.lastPlayedAt)}</Badge> : '';

  let header;

  if (showChannel) {
    header = (
      <div className="media">
        <div className="media-left media-middle">
          <Link to={channelUrl}>
            <Image
              className="media-object"
              src={image}
              errSrc="/static/podcast.png"
              imgProps={{
                height: 60,
                width: 60,
                alt: podcast.name,
              }}
            />
          </Link>
        </div>
        <div className="media-body">
          <h4>{showExpanded ? podcast.title :
            <Link to={podcastUrl}>{podcast.title}</Link>} {playedAt}</h4>
          <h5><Link to={channelUrl}>{podcast.name}</Link></h5>
        </div>
      </div>
    );
  } else {
    header = <h4><Link to={podcastUrl}>{podcast.title} {playedAt}</Link></h4>;
  }

  return (
    <Panel>
      {header}
      <div style={{ padding: 10 }}>
        <small>
          <time dateTime={podcast.pubDate}>{formatPubDate(podcast.pubDate)}</time>&nbsp;
          {podcast.source ? <a href={podcast.source} target="_blank">Source</a> : '' }
        </small>
        <Buttons {...props} />
      </div>
      {podcast.description && !showExpanded ?
      <Button
        className="form-control"
        title={podcast.isShowDetail ? 'Hide details' : 'Show details'}
        onClick={toggleDetail}
      ><Icon icon={podcast.isShowDetail ? 'chevron-up' : 'chevron-down'} />
      </Button> : ''}
    {podcast.description && (podcast.isShowDetail || showExpanded) ?
    <p
      className={showExpanded ? 'lead' : ''}
      style={{ marginTop: 20 }}
      dangerouslySetInnerHTML={sanitize(podcast.description)}
    /> : ''}
  </Panel>
  );
}


PodcastItem.propTypes = {
  podcast: PropTypes.object.isRequired,
  isLoggedIn: PropTypes.bool.isRequired,
  showChannel: PropTypes.bool.isRequired,
  showExpanded: PropTypes.bool,
  togglePlayer: PropTypes.func.isRequired,
  toggleDetail: PropTypes.func.isRequired,
  toggleBookmark: PropTypes.func.isRequired,
};
