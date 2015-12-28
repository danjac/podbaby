import React, { PropTypes } from 'react';
import { bindActionCreators } from 'react';
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

import { latest, player, bookmarks, subscribe } from '../actions';

import { sanitize } from './utils';

const ListItem = props => {
  const {
    podcast,
    createHref,
    isCurrentlyPlaying,
    setCurrentlyPlaying,
    unsubscribe,
    bookmark } = props;
  const url = createHref("/podcasts/channel/" + podcast.channelId + "/")
  // tbd get audio ref, set played at to last time
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
          <h4 className="media-heading"><a href={url}>{podcast.name}</a></h4>
          <Grid>
            <Row>
              <Col xs={6} md={6}>
                <h5>{podcast.title}</h5>
              </Col>
              <Col xs={6} mdPush={3} md={3}>
                <ButtonGroup>
                  <Button onClick={setCurrentlyPlaying}><Glyphicon glyph={ isCurrentlyPlaying ? 'stop': 'play' }  /></Button>
                  <a title="Download this podcast" className="btn btn-default" href={podcast.enclosureUrl}><Glyphicon glyph="download" /></a>
                  <Button onClick={bookmark} title={podcast.isBookmarked ? 'Remove bookmark' : 'Add to bookmarks'}>
                    <Glyphicon glyph={podcast.isBookmarked ? 'remove' : 'bookmark'} />
                  </Button>
                  <Button title="Unsubscribe from this channel" onClick={unsubscribe}><Glyphicon glyph="trash" /></Button>
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


export class Latest extends React.Component {

  componentDidMount() {
    const { dispatch } = this.props;
    dispatch(latest.getLatestPodcasts());
  }

  handleSelectPage(event, selectedEvent) {
    event.preventDefault();
    const { dispatch } = this.props;
    const page = selectedEvent.eventKey;
    dispatch(latest.getLatestPodcasts(page));
  }

  render() {
    const { dispatch } = this.props;
    const { createHref } = this.props.history;
    const { page, podcasts } = this.props;
    if (podcasts.length === 0) {
      return <div>You do not have any podcasts yet.</div>;
    }
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
        {podcasts.map(podcast => {
          const isCurrentlyPlaying = this.props.player.podcast && podcast.id === this.props.player.podcast.id;
          const setCurrentlyPlaying = (event) => {
            event.preventDefault();
            dispatch(player.setPodcast(isCurrentlyPlaying ? null : podcast));
          };
          const unsubscribe = (event) => {
            event.preventDefault();
            dispatch(subscribe.unsubscribe(podcast.channelId, podcast.name));
          };
          const bookmark = (event) => {
            event.preventDefault();
            const action = podcast.isBookmarked ? bookmarks.deleteBookmark : bookmarks.addBookmark;
            dispatch(action(podcast.id));
          };
          return <ListItem key={podcast.id}
                           podcast={podcast}
                           isCurrentlyPlaying={isCurrentlyPlaying}
                           unsubscribe={unsubscribe}
                           bookmark={bookmark}
                           setCurrentlyPlaying={setCurrentlyPlaying}
                           createHref={createHref} />;
        })}
        {pagination}
        </div>
    );
  }
}

Latest.propTypes = {
  podcasts: PropTypes.array.isRequired,
  page: PropTypes.object.isRequired,
  currentlyPlaying: PropTypes.number,
  dispatch: PropTypes.func.isRequired
};

const mapStateToProps = state => {
  const { podcasts, page } = state.latest;
  return {
    podcasts: podcasts || [],
    page: page,
    player: state.player
  };
};

export default connect(mapStateToProps)(Latest);
