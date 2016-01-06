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

  handleSelectPage(event, selectedEvent) {
    event.preventDefault();
    const page = selectedEvent.eventKey;
    this.actions.getBookmarks(page);
  }

  render() {
    return (
      <div>
        <form onSubmit={this.handleSearch.bind(this)}>
          <Input type="search"
                 ref="query"
                 placeholder="Find a podcast in your bookmarks" />
          <Button bsStyle="primary"
                  type="submit"
                  className="form-control"><Icon icon="search" /> Search</Button>
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
