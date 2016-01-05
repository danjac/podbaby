import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router';

import * as actions from '../actions';

import PodcastList from './podcasts';

export class Recent extends React.Component {

  componentDidMount() {
    const { dispatch } = this.props;
    dispatch(actions.plays.getRecentlyPlayed());
  }

  handleSelectPage(event, selectedEvent) {
    event.preventDefault();
    const { dispatch } = this.props;
    const page = selectedEvent.eventKey;
    dispatch(actions.plays.getRecentlyPlayed(page));
  }

  render() {

    return <PodcastList actions={actions}
                        ifEmpty="No recently played podcasts"
                        onSelectPage={this.handleSelectPage.bind(this)}
                        showChannel={true} {...this.props} />;
  }
}

Recent.propTypes = {
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

export default connect(mapStateToProps)(Recent);
