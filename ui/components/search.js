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
import PodcastList from './podcasts';
import { sanitize, formatPubDate } from './utils';

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

  componentDidMount() {
    const { q } = this.props.location.query;
    if (q) {
      this.props.dispatch(actions.search.search(q));
    }
  }

  componentWillReceiveProps(newProps) {
    const newQuery = newProps.location.query.q;
    const oldQuery = this.props.location.query.q;
    const isDiff = newQuery && oldQuery && newQuery !== oldQuery;
    if (isDiff) {
      console.log("fetching again", newQuery, oldQuery)
      this.props.dispatch(actions.search.search(newQuery));
    }
    return isDiff;
  }

  render() {

    const { dispatch, channels, podcasts, searchQuery } = this.props;
    const { createHref } = this.props.history;

    return (
      <div>
        {searchQuery ? <h2>Searching for "{searchQuery}"</h2> : ''}
        {channels.map(channel => {
          const subscribe = (event) => {
            event.preventDefault();
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
        {podcasts.length > 0 ? <hr /> : ''}
        <PodcastList actions={actions}
                     showChannel={true}
                     {...this.props} />
      </div>
    );
  }
}

const mapStateToProps = state => {
  const { podcasts, showDetail } = state.podcasts;
  const { query, channels } = state.search;
  return {
    searchQuery: query,
    podcasts: podcasts || [],
    channels,
    showDetail,
    player: state.player
  };
};

export default connect(mapStateToProps)(Search);
