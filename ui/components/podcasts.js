import _ from 'lodash';
import React, { PropTypes } from 'react';

import { Pagination } from 'react-bootstrap';

import Loading from './loading';
import Podcast from './podcast_item';

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
                          showExpanded={false}
                          toggleBookmark={toggleBookmark}
                          toggleDetail={toggleDetail}
                          togglePlayer={togglePlayer} />
        })}
        {pagination}
        </div>
      );
    }
}

export default PodcastList;
