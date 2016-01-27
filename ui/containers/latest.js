import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import DocumentTitle from 'react-document-title';
import { Button } from 'react-bootstrap';

import * as actions from '../actions';
import { podcastsSelector } from '../selectors';

import Icon from '../components/icon';
import PodcastList from '../components/podcasts';
import { getTitle } from './utils';

export class Latest extends React.Component {

  constructor(props) {
    super(props);
    this.actions = bindActionCreators(actions.latest, this.props.dispatch);
    this.handleSelectPage = this.handleSelectPage.bind(this);
    this.handleRefresh = this.handleRefresh.bind(this);
  }

  handleSelectPage(event, selectedEvent) {
    event.preventDefault();
    const page = selectedEvent.eventKey;
    this.actions.getLatestPodcasts(page);
  }

  handleRefresh(event) {
    event.prevetDefault();
    this.actions.getLatestPodcasts(1);
  }

  render() {
    const ifEmptyMsg = (
      <span>You haven't subscribed to any channels yet.
        Discover new channels and podcasts <Link to="/search/">here</Link>.</span>);

    return (
      <DocumentTitle title={getTitle('Latest podcasts')}>
        <div>
        {this.props.page.page !== 1 || this.props.isLoading ? '' :
        <form className="form">
          <Button bsStyle="primary" className="form-control" onClick={this.handleRefresh}>
            <Icon icon="refresh" /> Update
          </Button>
        </form>}
        <PodcastList
          actions={actions}
          ifEmpty={ifEmptyMsg}
          onSelectPage={this.handleSelectPage}
          showChannel {...this.props}
        />
        </div>
      </DocumentTitle>
    );
  }
}

Latest.propTypes = {
  podcasts: PropTypes.array.isRequired,
  page: PropTypes.object.isRequired,
  currentlyPlaying: PropTypes.number,
  dispatch: PropTypes.func.isRequired,
  isLoading: PropTypes.bool.isRequired,
  isLoggedIn: PropTypes.bool.isRequired,
};

const mapStateToProps = state => {
  const { page, isLoading } = state.podcasts;
  const { isLoggedIn } = state.auth;
  return {
    podcasts: podcastsSelector(state),
    isLoading,
    page,
    isLoggedIn,
  };
};

export default connect(mapStateToProps)(Latest);
