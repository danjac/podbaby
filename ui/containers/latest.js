import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import DocumentTitle from 'react-document-title';

import * as actions from '../actions';
import { podcastsSelector } from '../selectors';

import { getTitle } from '../components/utils';
import PodcastList from '../components/podcasts';

export class Latest extends React.Component {

  handleSelectPage(event, selectedEvent) {
    event.preventDefault();
    const { dispatch } = this.props;
    const page = selectedEvent.eventKey;
    dispatch(actions.latest.getLatestPodcasts(page));
  }

  render() {

    const ifEmptyMsg = (
      <span>You haven't subscribed to any channels yet.
        Discover new channels and podcasts <Link to="/podcasts/search/">here</Link>.</span>);

    return (
      <DocumentTitle title={getTitle('Latest podcasts')}>
        <PodcastList actions={actions}
                     ifEmpty={ifEmptyMsg}
                     onSelectPage={this.handleSelectPage.bind(this)}
                     showChannel={true} {...this.props} />
      </DocumentTitle>
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
  const { page, isLoading } = state.podcasts;
  const { isLoggedIn } = state.auth;
  return {
    podcasts: podcastsSelector(state),
    isLoading,
    page,
    isLoggedIn,
    player: state.player
  };
};

export default connect(mapStateToProps)(Latest);
