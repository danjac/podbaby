import _ from 'lodash';
import React, { PropTypes } from 'react';
import { Link } from 'react-router';

import {
  Grid,
  Row,
  Col,
  ButtonGroup,
  Button,
  Well,
  Panel,
  Pagination
} from 'react-bootstrap';

import Icon from './icon';
import Loading from './loading';
import { sanitize, formatPubDate } from './utils';

export class PodcastList extends React.Component {

  render() {
    const {
      actions,
      dispatch,
      podcasts,
      page,
      onSelectPage,
      player,
      isLoading,
      ifEmpty,
      showChannel
    } = this.props;


    if (isLoading) {
      return <Loading />;
    }

    if (_.isEmpty(podcasts)) {
      return <div>{ifEmpty || "No podcasts found"}</div>
    }

    const pagination = (
      page && onSelectPage && page.numPages > 1 ?
      <Pagination onSelect={onSelectPage}
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

          const isPlaying = player.podcast && podcast.id === player.podcast.id;

          const togglePlayer = event => {
            event.preventDefault();
            dispatch(actions.player.setPodcast(isPlaying ? null : podcast));
          };

          const toggleSubscribe = event => {
            event.preventDefault();
            const action = podcast.isSubscribed ? actions.subscribe.unsubscribe : actions.subscribe.subscribe;
            dispatch(action(podcast.channelId));
          };

          const toggleBookmark = event => {
            event.preventDefault();
            const action = podcast.isBookmarked ? actions.bookmarks.deleteBookmark : actions.bookmarks.addBookmark;
            dispatch(action(podcast.id));
          };

          const isShowingDetail = this.props.showDetail.includes(podcast.id);

          const toggleDetail = event => {
            event.preventDefault();
            const action = isShowingDetail ? actions.showDetail.hidePodcastDetail : actions.showDetail.showPodcastDetail;
            dispatch(action(podcast.id));
          };

          return <Podcast key={podcast.id}
                          podcast={podcast}
                          showChannel={showChannel}
                          toggleBookmark={toggleBookmark}
                          isShowingDetail={isShowingDetail}
                          toggleDetail={toggleDetail}
                          isPlaying={isPlaying}
                          togglePlayer={togglePlayer}
                          toggleSubscribe={toggleSubscribe} />
        })}
        {pagination}
        </div>
      );
    }
}

export const Podcast = props => {

  const {
    podcast,
    showChannel,
    isPlaying,
    isShowingDetail,
    togglePlayer,
    toggleSubscribe,
    toggleDetail,
    toggleBookmark } = props;

  const channelUrl = `/podcasts/channel/${podcast.channelId}/`;

  const header = showChannel ? <h3><Link to={channelUrl}>{podcast.name}</Link></h3> : <h3>{podcast.title}</h3>;

  return (
    <Panel bsStyle="primary" header={header}>
      <div className="media">
        {showChannel ?
        (<div className="media-left media-middle">
          <Link to={channelUrl}>
            <img className="media-object"
                 height={60}
                 width={60}
                 src={podcast.image}
                 alt={podcast.name} />
          </Link>
          </div> ) : ''}
        <div className="media-body">
          <Grid>
            <Row>
              <Col xs={6} md={6}>
                {showChannel ? <h4>{podcast.title}</h4> : ''}
                <br /><b>{formatPubDate(podcast.pubDate)}</b>
              </Col>
              <Col xs={6} mdPush={2} md={3}>
                <ButtonGroup>
                  <Button title={ isPlaying ? "Play": "Stop" } onClick={togglePlayer}><Icon icon={ isPlaying ? 'stop': 'play' }  /></Button>
                  <a title="Download this podcast" className="btn btn-default" href={podcast.enclosureUrl}><Icon icon="download" /></a>
                  <Button onClick={toggleBookmark} title={podcast.isBookmarked ? 'Remove bookmark' : 'Add to bookmarks'}>
                    <Icon icon={podcast.isBookmarked ? 'remove' : 'bookmark'} />
                  </Button>
                  {showChannel ? (<Button title={podcast.isSubscribed ? "Unsubscribe" : "Subscribe"} onClick={toggleSubscribe}>
                    <Icon icon={podcast.isSubscribed ? "trash" : "plus"} title={podcast.isSubscribed ? 'Unsubscribe' : 'Subscribe'} />
                  </Button>) : ''}
                </ButtonGroup>
              </Col>
            </Row>
          </Grid>
        </div>
      </div>
      {podcast.description ?
      (<div style={{paddingTop: "30px"}}>
        <Button className="form-control" onClick={toggleDetail}>
        <Icon icon={isShowingDetail ? 'chevron-up' : 'chevron-down'} />
        </Button>
      </div>) : ''}
      {podcast.description && isShowingDetail  ? <Well dangerouslySetInnerHTML={sanitize(podcast.description)} /> : ''}
  </Panel>
  );
};

export default PodcastList;
