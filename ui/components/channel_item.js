import React, { PropTypes } from 'react';
import { Link } from 'react-router';

import {
  Button,
  Panel,
} from 'react-bootstrap';

import Image from './image';
import Icon from './icon';
import { sanitize } from './utils';

function ChannelItem(props) {
  const { channel, subscribe, isLoggedIn } = props;
  const url = `/channel/${channel.id}/`;

  return (
    <Panel>
    <div className="media">
      <div className="media-left">
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
      <div className="media-body">
        <h4 className="media-heading">
          <Link to={url}>{channel.title}</Link>
        {isLoggedIn ?
          <Button style={{ float: 'right' }} title={channel.isSubscribed ?
            'Unsubscribe' : 'Subscribe'} onClick={subscribe}
          >
            <Icon icon={channel.isSubscribed ? 'unlink' : 'link'} /> {
            channel.isSubscribed ? 'Unsubscribe' : 'Subscribe'
            }
          </Button>
        : ''}
        </h4>
        <p
          style={{ marginTop: 20 }}
          dangerouslySetInnerHTML={sanitize(channel.description)}
        />
      </div>
    </div>
  </Panel>
  );
}

ChannelItem.propTypes = {
  channel: PropTypes.object.isRequired,
  subscribe: PropTypes.func.isRequired,
  isLoggedIn: PropTypes.bool.isRequired,
};

export default ChannelItem;
