import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router';

import * as actions from '../actions';

import PodcastList from './podcasts';

export class Latest extends React.Component {

  componentDidMount() {
    const { dispatch } = this.props;
    dispatch(actions.latest.getLatestPodcasts());
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
        Discover new channels and podcasts <Link to="/podcasts/search/">here</Link>.</span>);
    return <PodcastList actions={actions}
                        ifEmpty={ifEmptyMsg}
                        onSelectPage={this.handleSelectPage.bind(this)}
                        showChannel={true} {...this.props} />;
  }
}

Latest.propTypes = {
  podcasts: PropTypes.array.isRequired,
  page: PropTypes.object.isRequired,
  currentlyPlaying: PropTypes.number,
  dispatch: PropTypes.func.isRequired
};

const mapStateToProps = state => {
  const { podcasts, showDetail, page, isLoading } = state.podcasts;
  return {
    podcasts: podcasts || [],
    showDetail,
    isLoading,
    page,
    player: state.player
  };
};

export default connect(mapStateToProps)(Latest);
