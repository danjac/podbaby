import React from 'react';
import { connect } from 'react-redux';

import * as actions from '../actions';
import { bindAllActionCreators } from '../actions/utils';
import { Podcast } from './podcasts';
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
      isLoggedIn,
      player,
      bookmarks
    } = this.props;

    if (isLoading) {
      return <Loading />;
    }

    if (!podcast) {
      return <div>Sorry, no podcast found</div>;
    }

    // move these to a selector
    podcast.isBookmarked = bookmarks.includes(podcast.id);
    podcast.isPlaying = player.podcast && player.podcast.id === podcast.id;

    return <Podcast podcast={podcast}
                    showChannel={true}
                    showExpanded={true}
                    toggleBookmark={this.handleToggleBookmark.bind(this)}
                    toggleDetail={this.handleToggleDetail.bind(this)}
                    togglePlayer={this.handleTogglePlayer.bind(this)}
                    isLoggedIn={isLoggedIn} />
  }
}

const mapStateToProps = state => {

  const { podcast, isLoading } = state.podcast;
  const { isLoggedIn } = state.auth;
  const { bookmarks } = state.bookmarks;
  const player = state.player;

  return {
    podcast,
    isLoading,
    isLoggedIn,
    bookmarks,
    player
  };
};

export default connect(mapStateToProps)(PodcastDetail);
