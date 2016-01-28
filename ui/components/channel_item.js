import React, { PropTypes } from 'react';
import { Link } from 'react-router';

import {
  Button,
  Panel,
  Label,
} from 'react-bootstrap';

import Image from './image';
import Icon from './icon';
import { sanitize } from './utils';

function ChannelItem(props) {
  const { channel, subscribe, isLoggedIn } = props;
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
        <Link to={url}>
        <Image className="media-object"
          src={channel.image}
          errSrc="/static/podcast.png"
          imgProps={{
            height: 60,
            width: 60,
            alt: channel.title }}
        />
        </Link>
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
      dangerouslySetInnerHTML={sanitize(channel.description)}
    />
</Panel>
  );
}

ChannelItem.propTypes = {
  channel: PropTypes.object.isRequired,
  subscribe: PropTypes.func.isRequired,
  isLoggedIn: PropTypes.bool.isRequired,
};

export default ChannelItem;
