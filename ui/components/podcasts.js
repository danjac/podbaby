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

    const emptyMsg = typeof ifEmpty === "undefined" ? 'No podcasts found' : ifEmpty;
    if (_.isEmpty(podcasts)) {
      return <div>{emptyMsg}</div>
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
            dispatch(actions.bookmarks.toggleBookmark(podcast));
          };

          const isShowingDetail = this.props.showDetail.includes(podcast.id);

          const toggleDetail = event => {
            event.preventDefault();
            dispatch(actions.showDetail.toggleDetail(podcast, isShowingDetail));
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
  const title = <h5>{podcast.title}</h5>;

  let header = title;

  if (showChannel) {
    header = (
      <div>
        <h4><Link to={`/podcasts/channel/${podcast.channelId}/`}>{podcast.name}</Link></h4>
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
                  <Button title={ isPlaying ? "Stop": "Play" } onClick={togglePlayer}><Icon icon={ isPlaying ? 'stop': 'play' }  /></Button>
                  <a download
                     title="Download this podcast"
                     className="btn btn-default"
                     href={podcast.enclosureUrl}><Icon icon="download" /></a>
                  <Button onClick={toggleBookmark} title={podcast.isBookmarked ? 'Remove bookmark' : 'Add to bookmarks'}>
                    <Icon icon={podcast.isBookmarked ? 'bookmark' : 'bookmark-o'} />
                  </Button>
                </ButtonGroup>
              </Col>
            </Row>
          </Grid>
      </div>
      {podcast.description ?
      <Button className="form-control"
              title={isShowingDetail ? 'Hide details' : 'Show details'}
              onClick={toggleDetail}><Icon icon={isShowingDetail ? 'chevron-up': 'chevron-down'} /></Button> : ''}
    </div>
      {podcast.description && isShowingDetail  ? <Well style={{ marginTop: 20 }} dangerouslySetInnerHTML={sanitize(podcast.description)} /> : ''}
  </Panel>
  );
};

export default PodcastList;
