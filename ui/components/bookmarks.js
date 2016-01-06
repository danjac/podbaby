import _ from 'lodash';
import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

import { Button, Input } from 'react-bootstrap';
import Icon from './icon';

import * as actions from '../actions';
import PodcastList from './podcasts';


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
                            onSelectPage={this.handleSelectPage.bind(this)}
                            {...this.props} />
      </div>
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
  const { podcasts, page, showDetail, isLoading } = state.podcasts;
  return {
    podcasts: podcasts || [],
    showDetail,
    page,
    isLoading,
    query,
    player: state.player
  };
};

export default connect(mapStateToProps)(Bookmarks);
