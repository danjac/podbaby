import React, { PropTypes } from 'react';
import { bindActionCreators } from 'react';
import { connect } from 'react-redux';

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
    if (this.props.podcasts.length === 0) {
      return <div>You do not have any podcasts yet.</div>;
    }
    return <PodcastList actions={actions}
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
  const { podcasts, showDetail, page } = state.podcasts;
  return {
    podcasts: podcasts || [],
    showDetail: showDetail,
    page: page,
    player: state.player
  };
};

export default connect(mapStateToProps)(Latest);
