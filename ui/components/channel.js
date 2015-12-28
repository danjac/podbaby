import React, { PropTypes } from 'react';
import { connect } from 'react-redux';

import {
  Grid,
  Row,
  Col,
  Glyphicon,
  ButtonGroup,
  Button,
  Well,
  Pagination
} from 'react-bootstrap';

import * as actions from '../actions';
import { sanitize } from './utils';

const ListItem = props => {
  const {
    podcast,
    createHref,
    isCurrentlyPlaying,
    setCurrentlyPlaying,
    bookmark
  } = props;
  // tbd get audio ref, set played at to last time
  return (
    <div>
      <div className="media">
        <div className="media-body">
          <Grid>
            <Row>
              <Col xs={6} md={6}>
                <h4 className="media-heading">{podcast.title}</h4>
              </Col>
              <Col xs={6} mdPush={3} md={3}>
                <ButtonGroup>
                  <Button onClick={setCurrentlyPlaying}><Glyphicon glyph={ isCurrentlyPlaying ? 'stop': 'play' }  /> </Button>
                  <a className="btn btn-default" href={podcast.enclosureUrl}><Glyphicon glyph="download" /> </a>
                  <Button onClick={bookmark} title={podcast.isBookmarked ? 'Remove bookmark' : 'Add to bookmarks'}>
                    <Glyphicon glyph={podcast.isBookmarked ? 'remove' : 'bookmark'} />
                  </Button>
                </ButtonGroup>
              </Col>
            </Row>
          </Grid>
        </div>
      </div>
      {podcast.description ? <Well dangerouslySetInnerHTML={sanitize(podcast.description)} /> : ''}
    </div>
  );
};

export class Channel extends React.Component {

  componentDidMount(){
      this.props.dispatch(actions.channel.getChannel(this.props.params.id));
  }

  handleSubscribe(event) {
    event.preventDefault();
    const { channel, dispatch } = this.props;
    const action = channel.isSubscribed ? actions.subscribe.unsubscribe : actions.subscribe.subscribe;
    dispatch(action(channel.id, channel.title));
  }

  handleSelectPage(event, selectedEvent) {
    event.preventDefault();
    const { dispatch } = this.props;
    const page = selectedEvent.eventKey;
    dispatch(actions.channel.getChannel(this.props.params.id, page));
  }

  render() {
    const { channel, podcasts, page, dispatch, player } = this.props;
    if (!channel) {
      return <div></div>;
    }
    const isSubscribed = channel.isSubscribed;

    const pagination = (
      page.numPages > 1 ?
      <Pagination onSelect={this.handleSelectPage.bind(this)}
                  first
                  last
                  prev
                  next
                  maxButtons={6}
                  items={page.numPages}
                  activePage={page.page} /> : '');
    return (
      <div>
        {pagination}
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
            <Grid>
              <Row>
                <Col xs={6} md={9}>
                  <h2 className="media-heading">{channel.title}</h2>
                </Col>
                <Col xs={6} md={3}>
                  <ButtonGroup>
                    <Button title="Unsubscribe" onClick={this.handleSubscribe.bind(this)}><Glyphicon glyph="minus" /> Unsubscribe</Button>
                  </ButtonGroup>
                </Col>
              </Row>
            </Grid>
            {channel.description ? <Well dangerouslySetInnerHTML={sanitize(channel.description)} /> : ''}
          </div>
          </div>
          {podcasts.map(podcast => {
          const isCurrentlyPlaying = player.podcast && podcast.id === player.podcast.id;
          const setCurrentlyPlaying = event => {
            event.preventDefault();
            dispatch(actions.player.setPodcast(isCurrentlyPlaying ? null : podcast));
          };
          const bookmark = (event) => {
            event.preventDefault();
            const { bookmarks } = actions;
            const action = podcast.isBookmarked ? bookmarks.deleteBookmark : bookmarks.addBookmark;
            dispatch(action(podcast.id));
          };
          return <ListItem key={podcast.id}
                           podcast={podcast}
                           bookmark={bookmark}
                           isCurrentlyPlaying={isCurrentlyPlaying}
                           setCurrentlyPlaying={setCurrentlyPlaying}
                           channel={channel} />;
        })}
        {pagination}
      </div>
    );
  }
}

Channel.propTypes = {
  channel: PropTypes.object,
  podcasts: PropTypes.array,
  page: PropTypes.object,
  player: PropTypes.object,
  dispatch: PropTypes.func.isRequired
};

const mapStateToProps = state => {
  const { channel, podcasts, page } = state.channel;
  return {
    player: state.player,
    channel,
    podcasts,
    page
  };
};

export default connect(mapStateToProps)(Channel);
