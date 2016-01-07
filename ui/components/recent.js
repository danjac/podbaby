import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import { Button } from 'react-bootstrap';
import _ from 'lodash';

import * as actions from '../actions';
import { podcastsSelector } from '../selectors';

import Icon from './icon';
import PodcastList from './podcasts';

export class Recent extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindActionCreators(actions.plays, dispatch);
  }

  handleSelectPage(event, selectedEvent) {
    event.preventDefault();
    const page = selectedEvent.eventKey;
    this.actions.getRecentlyPlayed(page);
  }

  handleClearAll(event) {
    event.preventDefault();
    if (window.confirm("Are you sure you want to remove all the podcasts in your recently played list?")) {
      this.actions.clearAll();
    }
  }

  render() {

    return (
        <div>
          <PodcastList actions={actions}
                       ifEmpty="No recently played podcasts"
                       isLoggedIn={true}
                       onSelectPage={this.handleSelectPage.bind(this)}
                       showChannel={true} {...this.props} />
          {!_.isEmpty(this.props.podcasts) && !this.props.isLoading ?
          <Button className="form-control"
                  bsStyle="primary"
                  onClick={this.handleClearAll.bind(this)}>
                  <Icon icon="trash" /> Clear my recently played list
          </Button> : ''}
        </div>
      );

  }
}

Recent.propTypes = {
  podcasts: PropTypes.array.isRequired,
  page: PropTypes.object.isRequired,
  currentlyPlaying: PropTypes.number,
  dispatch: PropTypes.func.isRequired
};

const mapStateToProps = state => {
  const { page, isLoading } = state.podcasts;
  return {
    podcasts: podcastsSelector(state),
    isLoading,
    page,
  };
};

export default connect(mapStateToProps)(Recent);
