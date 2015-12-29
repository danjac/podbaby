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
  Pagination,
  Panel
} from 'react-bootstrap';

import { latest, player, bookmarks, subscribe, showDetail } from '../actions';

import { sanitize, formatPubDate } from './utils';

const ListItem = props => {
  const {
    podcast,
    createHref,
    isCurrentlyPlaying,
    setCurrentlyPlaying,
    unsubscribe,
    toggleDetail,
    isShowingDetail,
    bookmark } = props;
  const url = createHref("/podcasts/channel/" + podcast.channelId + "/")
  const header = <h3><a href={url}>{podcast.name}</a></h3>;

  return (
    <Panel header={header}>
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
                <h4>{podcast.title}</h4>
                <b>{formatPubDate(podcast.pubDate)}</b>
              </Col>
              <Col xs={6} mdPush={2} md={3}>
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
      <div style={{paddingTop: "30px"}}>
        <Button className="form-control" onClick={toggleDetail}>
        {isShowingDetail ? 'Show less' : 'Show more'} <Glyphicon glyph={isShowingDetail ? 'chevron-up' : 'chevron-down'} />
        </Button>
      </div>

      {podcast.description && isShowingDetail  ? <Well dangerouslySetInnerHTML={sanitize(podcast.description)} /> : ''}
  </Panel>
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

          const setCurrentlyPlaying = event => {
            event.preventDefault();
            dispatch(player.setPodcast(isCurrentlyPlaying ? null : podcast));
          };

          const unsubscribe = event => {
            event.preventDefault();
            dispatch(subscribe.unsubscribe(podcast.channelId, podcast.name));
          };

          const bookmark = event => {
            event.preventDefault();
            const action = podcast.isBookmarked ? bookmarks.deleteBookmark : bookmarks.addBookmark;
            dispatch(action(podcast.id));
          };

          const isShowingDetail = this.props.showDetail.includes(podcast.id);

          const toggleDetail = event => {
            event.preventDefault();
            const action = isShowingDetail ? showDetail.hidePodcastDetail : showDetail.showPodcastDetail;
            dispatch(action(podcast.id));
          };

          return <ListItem key={podcast.id}
                           podcast={podcast}
                           unsubscribe={unsubscribe}
                           bookmark={bookmark}
                           isShowingDetail={isShowingDetail}
                           toggleDetail={toggleDetail}
                           isCurrentlyPlaying={isCurrentlyPlaying}
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
  const { podcasts, showDetail, page } = state.podcasts;
  return {
    podcasts: podcasts || [],
    showDetail: showDetail,
    page: page,
    player: state.player
  };
};

export default connect(mapStateToProps)(Latest);
