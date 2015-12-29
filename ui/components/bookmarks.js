import React, { PropTypes } from 'react';
import { bindActionCreators } from 'react';
import { connect } from 'react-redux';


import * as actions from '../actions';
import PodcastList from './podcasts';


export class Bookmarks extends React.Component {

  componentDidMount() {
    const { dispatch } = this.props;
    dispatch(actions.bookmarks.getBookmarks());
  }

  handleSelectPage(event, selectedEvent) {
    event.preventDefault();
    const { dispatch } = this.props;
    const page = selectedEvent.eventKey;
    dispatch(actions.bookmarks.getBookmarks(page));
  }

  render() {
    const { page, podcasts, dispatch } = this.props;
    if (podcasts.length === 0) {
      return <div>You do not have any bookmarked podcasts yet.</div>;
    }

    return <PodcastList actions={actions}
                        showChannel={true} 
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
  const { podcasts, page, showDetail } = state.podcasts;
  return {
    podcasts: podcasts || [],
    showDetail,
    page: page,
    player: state.player
  };
};

export default connect(mapStateToProps)(Bookmarks);
