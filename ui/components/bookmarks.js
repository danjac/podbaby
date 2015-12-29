import React, { PropTypes } from 'react';
import { bindActionCreators } from 'react';
import { connect } from 'react-redux';


import * as actions from '../actions';
import PodcastList from './podcasts';


export class Bookmarks extends React.Component {

  componentDidMount() {
    const { dispatch } = this.props;
    dispatch(actions.bookmarks.getBookmarks());
    this.props.history.registerTransitionHook(this.handleLeavePage.bind(this));
  }

  handleLeavePage() {
    const { dispatch, isLoading } = this.props;
    if (!isLoading) {
      dispatch(actions.podcasts.unloadPodcasts());
    }
  }

  handleSelectPage(event, selectedEvent) {
    event.preventDefault();
    const { dispatch } = this.props;
    const page = selectedEvent.eventKey;
    dispatch(actions.bookmarks.getBookmarks(page));
  }

  render() {
    return <PodcastList actions={actions}
                        showChannel={true}
                        ifEmpty="You haven't added any bookmarks yet"
                        onSelectPage={this.handleSelectPage.bind(this)}
                        {...this.props} />;
  }
}

Bookmarks.propTypes = {
  podcasts: PropTypes.array.isRequired,
  page: PropTypes.object.isRequired,
  currentlyPlaying: PropTypes.number,
  dispatch: PropTypes.func.isRequired
};

const mapStateToProps = state => {
  const { podcasts, page, showDetail, isLoading } = state.podcasts;
  return {
    podcasts: podcasts || [],
    showDetail,
    page,
    isLoading,
    player: state.player
  };
};

export default connect(mapStateToProps)(Bookmarks);
