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
import { sanitize, formatPubDate } from './utils';

const PodcastItem = props => {
  const {
    podcast,
    createHref,
    isCurrentlyPlaying,
    setCurrentlyPlaying,
    bookmark,
    subscribe
  } = props;

  const url = createHref("/podcasts/channel/" + podcast.channelId + "/")
  return (
    <div>
      <div className="media">
        <div className="media-left media-middle">
          <a href={url}>
            <img className="media-object"
                 height={60}
                 width={60}
                 src={podcast.image}
                 alt={podcast.name} />
          </a>
        </div>
        <div className="media-body">
          <Grid>
            <Row>
              <Col xs={6} md={6}>
                <h4 className="media-heading"><a href={url}>{podcast.name}</a></h4>
              </Col>
              <Col xs={6} mdPush={3} md={3}>
                <b>{formatPubDate(podcast.pubDate)}</b><br />
                <ButtonGroup>
                  <Button onClick={setCurrentlyPlaying}><Glyphicon glyph={ isCurrentlyPlaying ? 'stop': 'play' }  /> </Button>
                  <a className="btn btn-default" href={podcast.enclosureUrl}><Glyphicon glyph="download" /> </a>
                  <Button onClick={bookmark} title={podcast.isBookmarked ? 'Remove bookmark' : 'Add to bookmarks'}>
                    <Glyphicon glyph={podcast.isBookmarked ? 'remove' : 'bookmark'} />
                  </Button>
                  <Button title={podcast.isSubscribed ? "Unsubscribe" : "Subscribe"} onClick={subscribe}>
                    <Glyphicon glyph={podcast.isSubscribed ? "trash" : "ok"} />
                  </Button>
                </ButtonGroup>
              </Col>
            </Row>
          </Grid>
        </div>
      </div>
      <h5>{podcast.title}</h5>
      {podcast.description ? <Well dangerouslySetInnerHTML={sanitize(podcast.description)} /> : ''}
    </div>
  );
};

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
    const { q } = newProps.location.query;
    const isDiff = this.props.searchQuery !== q;
    if (isDiff) {
      this.props.dispatch(actions.search.search(q));
    }
    return isDiff;
  }

  render() {

    const { dispatch, channels, podcasts, player, searchQuery } = this.props;
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
        {podcasts.map(podcast => {

          const subscribe = event => {
            event.preventDefault();
            const action = podcast.isSubscribed ? actions.subscribe.unsubscribe : actions.subscribe.subscribe;
            dispatch(action(podcast.channelId, podcast.name));
          };

          const isCurrentlyPlaying = player.podcast && podcast.id === player.podcast.id;

          const setCurrentlyPlaying = event => {
            event.preventDefault();
            dispatch(actions.player.setPodcast(isCurrentlyPlaying ? null : podcast));
          };

          const bookmark = event => {
            event.preventDefault();
            const { bookmarks } = actions;
            const action = podcast.isBookmarked ? bookmarks.deleteBookmark : bookmarks.addBookmark;
            dispatch(action(podcast.id));
          };

          return (
            <PodcastItem key={podcast.id}
                         podcast={podcast}
                         subscribe={subscribe}
                         bookmark={bookmark}
                         isCurrentlyPlaying={isCurrentlyPlaying}
                         setCurrentlyPlaying={setCurrentlyPlaying}
                         createHref={createHref} />
          );
        })}
      </div>
    );
  }
}

const mapStateToProps = state => {
  const { podcasts } = state.podcasts;
  const { query, channels } = state.search;
  return {
    searchQuery: query,
    podcasts,
    channels,
    player: state.player
  };
};

export default connect(mapStateToProps)(Search);
