import React from 'react';
import { connect } from 'react-redux';
import DocumentTitle from 'react-document-title';

import * as actions from '../actions';
import { podcastSelector } from '../selectors';
import { bindAllActionCreators } from '../actions/utils';
import { Podcast } from './podcasts';
import { getTitle } from './utils';
import Loading from './loading';

class PodcastDetail extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindAllActionCreators(actions, dispatch);
  }

  handleTogglePlayer(event) {
    event.preventDefault();
    this.actions.player.togglePlayer(this.props.podcast);
  }

  handleToggleBookmark(event) {
    event.preventDefault();
    this.actions.bookmarks.toggleBookmark(this.props.podcast)
  }

  handleToggleDetail(event) {
    event.preventDefault();
    this.actions.showDetail.toggleDetail(this.props.podcast)
  }

  render() {
    const {
      podcast,
      isLoading,
      isLoggedIn
    } = this.props;

    if (isLoading) {
      return <Loading />;
    }

    if (!podcast) {
      return <div>Sorry, no podcast found</div>;
    }

    return (
      <DocumentTitle title={getTitle(podcast.name, podcast.title)}>
      <Podcast podcast={podcast}
               showChannel={true}
               showExpanded={true}
               toggleBookmark={this.handleToggleBookmark.bind(this)}
               toggleDetail={this.handleToggleDetail.bind(this)}
               togglePlayer={this.handleTogglePlayer.bind(this)}
               isLoggedIn={isLoggedIn} />
      </DocumentTitle>
    );
  }
}

const mapStateToProps = state => {

  const podcast = podcastSelector(state);

  const { isLoading } = state.podcast;
  const { isLoggedIn } = state.auth;
  const { bookmarks } = state.bookmarks;

  return {
    podcast,
    isLoading,
    isLoggedIn
  };
};

export default connect(mapStateToProps)(PodcastDetail);
