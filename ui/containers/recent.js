import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Button } from 'react-bootstrap';
import DocumentTitle from 'react-document-title';

import * as actions from '../actions';
import { podcastsSelector } from '../selectors';

import Icon from '../components/icon';
import PodcastList from '../components/podcasts';
import { getTitle } from './utils';

export class Recent extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindActionCreators(actions.plays, dispatch);
    this.handleSelectPage = this.handleSelectPage.bind(this);
    this.handleClearAll = this.handleClearAll.bind(this);
  }

  handleSelectPage(event, selectedEvent) {
    event.preventDefault();
    const page = selectedEvent.eventKey;
    this.actions.getRecentlyPlayed(page);
  }

  handleClearAll(event) {
    event.preventDefault();
    if (window.confirm(
      'Are you sure you want to remove all the podcasts in your recently played list?')) {
      this.actions.clearAll();
    }
  }

  render() {
    return (
      <DocumentTitle title={getTitle('My recently played podcasts')}>
        <div>
          <PodcastList
            actions={actions}
            ifEmpty="No recently played podcasts"
            isLoggedIn
            onSelectPage={this.handleSelectPage}
            showChannel
            {...this.props}
          />
          {this.props.podcasts.size > 0 && !this.props.isLoading ?
          <Button
            className="form-control"
            bsStyle="primary"
            onClick={this.handleClearAll}
          >
            <Icon icon="trash" /> Clear my recently played list
          </Button> : ''}
        </div>
      </DocumentTitle>
      );
  }
}

Recent.propTypes = {
  podcasts: PropTypes.object.isRequired,
  page: PropTypes.object.isRequired,
  currentlyPlaying: PropTypes.number,
  dispatch: PropTypes.func.isRequired,
  isLoading: PropTypes.bool.isRequired,
};

const mapStateToProps = state => {
  return {
    podcasts: podcastsSelector(state),
    isLoading: state.podcasts.get('isLoading'),
    page: state.podcasts.get('page'),
  };
};

export default connect(mapStateToProps)(Recent);
