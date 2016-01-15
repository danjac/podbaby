import React, { PropTypes } from 'react';
import { Link } from 'react-router';

import {
  Grid,
  Row,
  Col,
  ButtonGroup,
  Button,
  Panel,
} from 'react-bootstrap';

import Image from './image';
import Icon from './icon';

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
        <Grid>
          <Row>
            <Col xs={6} md={9}>
              <h4 className="media-heading"><Link to={url}>{channel.title}</Link></h4>
            </Col>
            <Col xs={6} md={3}>
              {isLoggedIn ?
              <ButtonGroup>
                <Button title={channel.isSubscribed ?
                  'Unsubscribe' : 'Subscribe'} onClick={subscribe}
                >
                  <Icon icon={channel.isSubscribed ? 'unlink' : 'link'} /> {
                  channel.isSubscribed ? 'Unsubscribe' : 'Subscribe'
                  }
                </Button>
              </ButtonGroup>
              : ''}
            </Col>
          </Row>
        </Grid>
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
