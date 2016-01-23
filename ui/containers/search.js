import _ from 'lodash';
import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Button, Input, Tabs, Tab } from 'react-bootstrap';
import { routeActions } from 'redux-simple-router';
import DocumentTitle from 'react-document-title';

import * as actions from '../actions';
import { bindAllActionCreators } from '../actions/utils';
import { podcastsSelector, channelsSelector } from '../selectors';
import ChannelItem from '../components/channel_item';
import PodcastList from '../components/podcasts';
import Icon from '../components/icon';
import { getTitle } from './utils';

export class Search extends React.Component {

  constructor(props) {
    super(props);
    this.actions = bindAllActionCreators(actions, this.props.dispatch);
    this.route = bindActionCreators(routeActions, this.props.dispatch);
    this.handleSearch = this.handleSearch.bind(this);
    this.handleSelect = this.handleSelect.bind(this);
  }

  handleSearch(event) {
    event.preventDefault();
    const value = _.trim(this.refs.query.getValue());
    if (value && value !== this.props.searchQuery) {
      this.route.replace(`/search/?q=${value}`);
      this.actions.search.search(value);
    } else {
      this.actions.search.clearSearch();
    }
  }

  handleSelect() {
    this.refs.query.getInputDOMNode().select();
  }

  renderSearchResults() {
    const {
      dispatch,
      channels,
      podcasts,
      isLoading,
      searchQuery } = this.props;

    if (isLoading) {
      return '';
    }

    if (channels.length === 0 &&
        podcasts.length === 0 &&
        searchQuery) return <div>Sorry, no results found for your search.</div>;

    const channelItems = channels.length > 0 && channels.map(channel => {
      const subscribe = event => {
        event.preventDefault();
        dispatch(actions.subscribe.toggleSubscribe(channel));
      };

      return (
        <ChannelItem
          key={channel.id}
          channel={channel}
          subscribe={subscribe}
          {...this.props}
        />
      );
    });

    const podcastItems = podcasts.length > 0 ?

      <PodcastList
        actions={actions}
        showChannel
        ifEmpty=""
        {...this.props}
      /> : '';

    if (podcastItems && channelItems) {
      const tabStyle = { marginTop: 20 };

      return (
        <Tabs defaultActiveKey={1}>
          <Tab eventKey={1} title="Podcasts" style={tabStyle}>
            {podcastItems}
          </Tab>
          <Tab eventKey={2} title="Feeds" style={tabStyle}>
            {channelItems}
          </Tab>
        </Tabs>
      );
    } else if (channelItems) {
      return <div>{channelItems}</div>;
    } else if (podcastItems) {
      return podcastItems;
    }
  }

  render() {
    const { searchQuery } = this.props;

    const help = (
      searchQuery ? '' :
        <span>
          <b>Hint:</b>
          Try a general category e.g. <em>history</em> or <em>movies</em>,
          the title of a podcast, or the name of a feed e. g. <em>RadioLab</em>.
        </span>
      );

    return (
      <DocumentTitle title={getTitle('Search podcasts and feeds')}>
        <div>
          <form className="form" onSubmit={this.handleSearch}>
            <Input
              type="search"
              ref="query"
              defaultValue={searchQuery}
              help={help}
              onClick={this.handleSelect}
              placeholder="Find a feed or podcast"
            />
            <Button type="submit" bsStyle="primary" className="form-control">
              <Icon icon="search" /> Search
            </Button>
          </form>
          {this.renderSearchResults()}
        </div>
      </DocumentTitle>
    );
  }

}

Search.propTypes = {
  dispatch: PropTypes.func.isRequired,
  location: PropTypes.object.isRequired,
  channels: PropTypes.array.isRequired,
  podcasts: PropTypes.array.isRequired,
  isLoading: PropTypes.bool.isRequired,
  searchQuery: PropTypes.string.isRequired,
};

const mapStateToProps = state => {
  const { isLoading } = state.podcasts;
  const { query } = state.search;
  const { isLoggedIn } = state.auth;

  const podcasts = podcastsSelector(state);
  const { channels } = channelsSelector(state);

  return {
    searchQuery: query,
    podcasts,
    channels,
    isLoading,
    isLoggedIn,
  };
};

export default connect(mapStateToProps)(Search);
