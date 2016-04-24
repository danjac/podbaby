import React, { PropTypes } from 'react';
import { connect } from 'react-redux';
import DocumentTitle from 'react-document-title';

import * as actions from '../actions';
import { podcastSelector } from '../selectors';
import { bindAllActionCreators } from '../actions/utils';
import Podcast from '../components/podcast_item';
import Loading from '../components/loading';
import { getTitle } from './utils';

class PodcastDetail extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindAllActionCreators(actions, dispatch);
    this.handleToggleDetail = this.handleToggleDetail.bind(this);
    this.handleToggleBookmark = this.handleToggleBookmark.bind(this);
    this.handleTogglePlayer = this.handleTogglePlayer.bind(this);
  }

  handleTogglePlayer(event) {
    event.preventDefault();
    this.actions.player.togglePlayer(this.props.podcast);
  }

  handleToggleBookmark(event) {
    event.preventDefault();
    this.actions.bookmarks.toggleBookmark(this.props.podcast);
  }

  handleToggleDetail(event) {
    event.preventDefault();
    this.actions.showDetail.toggleDetail(this.props.podcast);
  }

  render() {
    const {
      podcast,
      isLoading,
      isLoggedIn,
    } = this.props;

    if (isLoading) {
      return <Loading />;
    }

    if (!podcast) {
      return <div>Sorry, no podcast found</div>;
    }

    return (
      <DocumentTitle title={getTitle(podcast.name, podcast.title)}>
      <Podcast
        podcast={podcast}
        showChannel
        showExpanded
        showImage
        toggleBookmark={this.handleToggleBookmark}
        toggleDetail={this.handleToggleDetail}
        togglePlayer={this.handleTogglePlayer}
        isLoggedIn={isLoggedIn}
      />
      </DocumentTitle>
    );
  }
}

PodcastDetail.propTypes = {
  dispatch: PropTypes.func.isRequired,
  podcast: PropTypes.object.isRequired,
  isLoading: PropTypes.bool.isRequired,
  isLoggedIn: PropTypes.bool.isRequired,
};

const mapStateToProps = state => {
  const podcast = podcastSelector(state);

  const { isLoading } = state.podcast;
  const { isLoggedIn } = state.auth;

  return {
    podcast,
    isLoading,
    isLoggedIn,
  };
};

export default connect(mapStateToProps)(PodcastDetail);
