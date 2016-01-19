import React, { PropTypes } from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import DocumentTitle from 'react-document-title';

import * as actions from '../actions';
import { podcastsSelector } from '../selectors';

import PodcastList from '../components/podcasts';
import { getTitle } from './utils';

export class Latest extends React.Component {

  constructor(props) {
    super(props);
    this.handleSelectPage = this.handleSelectPage.bind(this);
  }

  handleSelectPage(event, selectedEvent) {
    event.preventDefault();
    const { dispatch } = this.props;
    const page = selectedEvent.eventKey;
    dispatch(actions.latest.getLatestPodcasts(page));
  }

  render() {
    const ifEmptyMsg = (
      <span>You haven't subscribed to any channels yet.
        Discover new channels and podcasts <Link to="/search/">here</Link>.</span>);

    return (
      <DocumentTitle title={getTitle('Latest podcasts')}>
        <PodcastList
          actions={actions}
          ifEmpty={ifEmptyMsg}
          onSelectPage={this.handleSelectPage}
          showChannel {...this.props}
        />
      </DocumentTitle>
    );
  }
}

Latest.propTypes = {
  podcasts: PropTypes.object.isRequired,
  page: PropTypes.object.isRequired,
  currentlyPlaying: PropTypes.number,
  dispatch: PropTypes.func.isRequired,
  isLoading: PropTypes.bool.isRequired,
  isLoggedIn: PropTypes.bool.isRequired,
};

const mapStateToProps = state => {
  return {
    podcasts: podcastsSelector(state),
    page: state.podcasts.get('page'),
    isLoading: state.podcasts.get('isLoading'),
    isLoggedIn: state.auth.get('isLoggedIn'),
  };
};

export default connect(mapStateToProps)(Latest);
