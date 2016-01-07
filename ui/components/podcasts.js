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
      isLoggedIn,
      podcasts,
      page,
      onSelectPage,
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

          const togglePlayer = event => {
            event.preventDefault();
            dispatch(actions.player.togglePlayer(podcast));
          };

          const toggleBookmark = event => {
            event.preventDefault();
            dispatch(actions.bookmarks.toggleBookmark(podcast));
          };

          const toggleDetail = event => {
            event.preventDefault();
            dispatch(actions.showDetail.toggleDetail(podcast));
          };

          return <Podcast key={podcast.id}
                          isLoggedIn={isLoggedIn}
                          podcast={podcast}
                          showChannel={showChannel}
                          toggleBookmark={toggleBookmark}
                          toggleDetail={toggleDetail}
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
    showExpanded,
    isLoggedIn,
    togglePlayer,
    toggleDetail,
    toggleBookmark } = props;

  const channelUrl = `/channel/${podcast.channelId}/`;
  const podcastUrl = `/podcast/${podcast.id}/`;

  let header;

  if (showChannel) {
    header = (
      <div>
        <h4><Link to={showExpanded ? channelUrl : podcastUrl}>{podcast.name}</Link></h4>
        <h5>{podcast.title}</h5>
      </div>
    );
  } else {
    header = <h5><Link to={podcastUrl}>{podcast.title}</Link></h5>;
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
                  <Button title={ podcast.isPlaying ? "Stop": "Play" } onClick={togglePlayer}><Icon icon={ podcast.isPlaying ? 'stop': 'play' }  /></Button>
                  <a download
                     title="Download this podcast"
                     className="btn btn-default"
                     href={podcast.enclosureUrl}><Icon icon="download" /></a>
                   {isLoggedIn ? <Button onClick={toggleBookmark} title={podcast.isBookmarked ? 'Remove bookmark' : 'Add to bookmarks'}>
                    <Icon icon={podcast.isBookmarked ? 'bookmark' : 'bookmark-o'} />
                  </Button> : ''}
                </ButtonGroup>
              </Col>
            </Row>
          </Grid>
      </div>
      {podcast.description && !showExpanded ?
      <Button className="form-control"
              title={podcast.isShowDetail ? 'Hide details' : 'Show details'}
              onClick={toggleDetail}><Icon icon={podcast.isShowDetail ? 'chevron-up': 'chevron-down'} /></Button> : ''}
    </div>
      {podcast.description && (podcast.isShowDetail || showExpanded)  ? <Well style={{ marginTop: 20 }} dangerouslySetInnerHTML={sanitize(podcast.description)} /> : ''}
  </Panel>
  );
};

export default PodcastList;
