import React, { PropTypes } from 'react';
import { Link } from 'react-router';

import {
  Button,
  Panel,
  Label,
} from 'react-bootstrap';

import Image from './image';
import Icon from './icon';
import { sanitize, highlight } from './utils';

function ChannelItem(props) {
  const { channel, subscribe, isLoggedIn, searchQuery } = props;
  const url = `/channel/${channel.id}/`;

  return (
    <Panel>
    <div className="thumbnail">
        <div className="caption text-center">
          <h4>
            <Link to={url}>{channel.title}</Link>
              </h4>
          {channel.numPodcasts ?
          <h5>
            <Label bsStyle="primary">
              {channel.numPodcasts} podcast{channel.numPodcasts > 1 ? 's' : ''}
            </Label>
          </h5> : ''}
        </div>
        {props.showImage ?
        <Link to={url}>
          <Image
            hideIfMobile
            className="media-object"
            src={channel.image}
            errSrc="/static/podcast.png"
            imgProps={{
              height: 60,
              width: 60,
              alt: channel.title }}
          />
        </Link> : ''}
      </div>
    {isLoggedIn ?
    <div style={{ marginTop: 20 }}>
      <Button
        bsStyle={channel.isSubscribed ? 'default' : 'primary'}
        className="form-control"
        title={channel.isSubscribed ?
        'Unsubscribe' : 'Subscribe'} onClick={subscribe}
      >
        <Icon icon={channel.isSubscribed ? 'unlink' : 'link'} /> {
        channel.isSubscribed ? 'Unsubscribe' : 'Subscribe'
        }
      </Button>
    </div>
    : ''}
    <p
      style={{ marginTop: 20 }}
      dangerouslySetInnerHTML={sanitize(highlight(channel.description, searchQuery))}
    />
</Panel>
  );
}

ChannelItem.propTypes = {
  channel: PropTypes.object.isRequired,
  showImage: PropTypes.bool.isRequired,
  subscribe: PropTypes.func.isRequired,
  isLoggedIn: PropTypes.bool.isRequired,
  searchQuery: PropTypes.string,
};

export default ChannelItem;
