import _ from 'lodash';
import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import DocumentTitle from 'react-document-title';

import { Button, Input } from 'react-bootstrap';

import * as actions from '../actions';
import { podcastsSelector } from '../selectors';
import PodcastList from './podcasts';
import Icon from './icon';
import { getTitle } from './utils';


export class Bookmarks extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindActionCreators(actions.bookmarks, dispatch);
  }

  handleSearch(event) {
    event.preventDefault();
    const query = _.trim(this.refs.query.getValue());
    if (query) {
      this.actions.searchBookmarks(query);
    } else {
      this.actions.getBookmarks();
    }
  }

  handleClearSearch(event) {
    event.preventDefault();
    this.refs.query.getInputDOMNode().value = "";
    this.actions.getBookmarks();
  }

  handleClickSearch(event) {
    event.preventDefault();
    this.refs.query.getInputDOMNode().select();
  }

  handleSelectPage(event, selectedEvent) {
    event.preventDefault();
    const page = selectedEvent.eventKey;
    this.actions.getBookmarks(page);
  }

  render() {
    const { query } = this.props;
    return (
      <DocumentTitle title={getTitle('My bookmarks')}>
      <div>
        <form onSubmit={this.handleSearch.bind(this)}>
          <Input type="search"
                 ref="query"
                 onClick={this.handleClickSearch.bind(this)}
                 placeholder="Find a podcast in your bookmarks" />
          <Input>
            <Button bsStyle="primary"
                    type="submit"
                    defaultValue={query}
                    className="form-control"><Icon icon="search" /> Search</Button>
          </Input>
          {query ? <Input>
            <Button bsStyle="default"
                    onClick={this.handleClearSearch.bind(this)}
                    className="form-control"><Icon icon="refresh" /> Show all bookmarks</Button>
          </Input> : ''}
        </form>
        <PodcastList actions={actions}
                            showChannel={true}
                            ifEmpty="No bookmarks found"
                            isLoggedIn={true}
                            onSelectPage={this.handleSelectPage.bind(this)}
                            {...this.props} />
      </div>
    </DocumentTitle>
    );
  }
}

Bookmarks.propTypes = {
  podcasts: PropTypes.array.isRequired,
  page: PropTypes.object.isRequired,
  currentlyPlaying: PropTypes.number,
  dispatch: PropTypes.func.isRequired
};

const mapStateToProps = state => {
  const { query } = state.bookmarks;
  const { page, isLoading } = state.podcasts;
  return {
    podcasts: podcastsSelector(state),
    page,
    isLoading,
    query,
  };
};

export default connect(mapStateToProps)(Bookmarks);
