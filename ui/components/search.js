import React, { PropTypes } from 'react';
import { connect } from 'react-redux';

import {
  Grid,
  Row,
  Col,
  Glyphicon,
  ButtonGroup,
  Button,
  Well
} from 'react-bootstrap';


import  * as actions from '../actions';

const ChannelItem = props => {
  const { channel, createHref, subscribe } = props;
  return (
    <div className="media">
      <div className="media-left">
        <a href="#">
          <img className="media-object"
               height={60}
               width={60}
               src={channel.image}
               alt={channel.title} />
        </a>
      </div>
      <div className="media-body">
        <h4 className="media-heading"><a href={createHref("/podcasts/channel/" + channel.id + "/")}>{channel.title}</a></h4>
        <Grid>
          <Row>
            <Col xs={6} md={9}>
              <Well>{channel.description}</Well>
            </Col>
            <Col xs={6} md={3}>
              <ButtonGroup>
                <Button title={channel.isSubscribed ? "Unsubscribe" : "Subscribe"} onClick={subscribe}>
                  <Glyphicon glyph={channel.isSubscribed ? "trash" : "ok"} /> {channel.isSubscribed ? 'Unsubscribe' : 'Subscribe'}
                </Button>
              </ButtonGroup>
            </Col>
          </Row>
        </Grid>
      </div>
    </div>
  );
};


export class Search extends React.Component {
  render() {

    const { dispatch, channels, searchQuery } = this.props;
    const { createHref } = this.props.history;
    console.log(searchQuery, channels)

    return (
      <div>
        <h2>Searching for {searchQuery}</h2>
        {channels.map(channel => {
          const subscribe = () => {
              const action = channel.isSubscribed ? actions.subscribe.unsubscribe : actions.subscribe.subscribe;
              dispatch(action(channel.id, channel.title));
          };
          return (
            <ChannelItem key={channel.id}
                         channel={channel}
                         subscribe={subscribe}
                         createHref={createHref} />
          );
        })}
      </div>
    );
  }
}

const mapStateToProps = state => {
  const { query, podcasts, channels, numResults } = state.search;
  return {
    searchQuery: query,
    podcasts,
    channels,
    numResults
  };
};

export default connect(mapStateToProps)(Search);
