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
import Image from './image';
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
            dispatch(actions.player.setPodcast(player, isPlaying ? null : podcast));
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
                          togglePlayer={togglePlayer} />
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
    toggleDetail,
    toggleBookmark } = props;

  const channelUrl = `/podcasts/channel/${podcast.channelId}/`;
  const title = <h4>{podcast.title}</h4>;

  let header = title;

  if (showChannel) {
    header = (
      <div>
        <h3><Link to={`/podcasts/channel/${podcast.channelId}/`}>{podcast.name}</Link></h3>
        {title}
      </div>
    );
  }

  return (
    <Panel>
      <div className="media">
        {showChannel ? (
        <div className="media-left media-middle">
          <Link to={channelUrl}>
              <Image className="media-object"
                     src={podcast.image}
                     errSrc='/static/podcast.png'
                     imgProps={{
                     height:60,
                     width:60,
                     alt:podcast.name }} />
          </Link>
          </div>
          ) : '' }
        <div className="media-body">
          <Grid>
            <Row>
              <Col xs={6} md={9}>
              {header}
              <p><small><time dateTime={podcast.pubDate}>{formatPubDate(podcast.pubDate)}</time></small></p>
              </Col>
              <Col xs={6} md={3}>
                <ButtonGroup>
                  <Button title={ isPlaying ? "Play": "Stop" } onClick={togglePlayer}><Icon icon={ isPlaying ? 'stop': 'play' }  /></Button>
                  <a title="Download this podcast" className="btn btn-default" href={podcast.enclosureUrl}><Icon icon="download" /></a>
                  <Button onClick={toggleBookmark} title={podcast.isBookmarked ? 'Remove bookmark' : 'Add to bookmarks'}>
                    <Icon icon={podcast.isBookmarked ? 'remove' : 'bookmark'} />
                  </Button>
                  {podcast.description ?
                  <Button title={isShowingDetail ? 'Hide details' : 'Show details'} onClick={toggleDetail}><Icon icon={isShowingDetail ? 'compress': 'expand'} /></Button>
                  : ''}
                </ButtonGroup>
              </Col>
            </Row>
          </Grid>
        </div>
      </div>
      {podcast.description && isShowingDetail  ? <Well style={{ marginTop: 20 }} dangerouslySetInnerHTML={sanitize(podcast.description)} /> : ''}
  </Panel>
  );
};

export default PodcastList;
