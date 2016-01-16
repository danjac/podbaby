import React, { PropTypes } from 'react';
import { Link } from 'react-router';
import { sanitize, formatPubDate } from './utils';

import {
  Grid,
  Row,
  Col,
  ButtonGroup,
  Button,
  Panel,
} from 'react-bootstrap';


import Icon from './icon';
import Image from './image';


const Buttons = props => {
  const { podcast } = props;
  return (
    <ButtonGroup vertical={props.vertical}>
     <Button
       title={ podcast.isPlaying ? 'Stop' : 'Play' }
       onClick={props.togglePlayer}
     ><Icon icon={ podcast.isPlaying ? 'stop' : 'play' } />
     </Button>
     <a
       download
       title="Download this podcast"
       className="btn btn-default"
       href={podcast.enclosureUrl}
     ><Icon icon="download" /></a>
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
  } = props;

  const channelUrl = `/channel/${podcast.channelId}/`;
  const podcastUrl = `/podcast/${podcast.id}/`;
  const image = podcast.image || '/static/podcast.png';

  let header;

  if (showChannel) {
    header = (
      <div>
        <h4>{showExpanded ? podcast.title : <Link to={podcastUrl}>{podcast.title}</Link>}</h4>
        <h5><Link to={channelUrl}>{podcast.name}</Link></h5>
      </div>
    );
  } else {
    header = <h4><Link to={podcastUrl}>{podcast.title}</Link></h4>;
  }

  return (
    <Panel>
      <div className="media">
        {showChannel ? (
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
          ) : '' }
        <div className="media-body">
          <Grid>
            <Row>
              <Col xs={6} md={9}>
              {header}
              <p><small>
                <time dateTime={podcast.pubDate}>{formatPubDate(podcast.pubDate)}</time>&nbsp;
                {podcast.source ? <a href={podcast.source} target="_blank">Source</a> : '' }
              </small></p>
              </Col>
              <Col className="hidden-xs hidden-sm" md={3}>
                <Buttons {...props} />
              </Col>
              <Col className="hidden-md hidden-lg" xs={3} sm={3}>
                <Buttons {...props} vertical />
              </Col>
            </Row>
          </Grid>
      </div>
      {podcast.description && !showExpanded ?
      <Button
        className="form-control"
        title={podcast.isShowDetail ? 'Hide details' : 'Show details'}
        onClick={toggleDetail}
      ><Icon icon={podcast.isShowDetail ? 'chevron-up' : 'chevron-down'} />
      </Button> : ''}
    </div>
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
