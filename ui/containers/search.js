import _ from 'lodash';
import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Button, Input, Tabs, Tab } from 'react-bootstrap';
import DocumentTitle from 'react-document-title';

import  * as actions from '../actions';
import { podcastsSelector, channelsSelector } from '../selectors';
import ChannelItem from '../components/channel_item';
import PodcastList from '../components/podcasts';
import Icon from '../components/icon';
import { getTitle } from './utils';

export class Search extends React.Component {

  constructor(props) {
    super(props);
    const { search } = bindActionCreators(actions.search, this.props.dispatch);
    this.search = search;
  }

  componentDidMount() {
    const query = this.props.location.query.q || "";
    this.search(query);
    this.refs.query.getInputDOMNode().focus();
  }

  handleSearch(event) {
    event.preventDefault();
    const value = this.refs.query.getValue();
    this.search(_.trim(value));
  }

  handleFocus(event) {
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

      if (channels.length == 0 &&
          podcasts.length == 0 &&
          searchQuery) return <div>Sorry, no results found for your search.</div>;

      const channelItems = channels.length > 0 && channels.map(channel => {

        const subscribe = event => {
          event.preventDefault();
          dispatch(actions.subscribe.toggleSubscribe(channel));
        };

        return (
          <ChannelItem key={channel.id}
                       channel={channel}
                       subscribe={subscribe}
                       {...this.props} />
        );

      });

      const podcastItems = podcasts.length > 0 ?

        <PodcastList actions={actions}
                     showChannel={true}
                     ifEmpty=''
                      {...this.props} /> : '';

      if (podcastItems && channelItems) {

        const tabStyle = { marginTop: 20 };

        return (
          <Tabs defaultActiveKey={1}>
            <Tab eventKey={1} title="Podcasts" style={tabStyle}>
              {podcastItems}
            </Tab>
            <Tab eventKey={2} title="Channels" style={tabStyle}>
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
          <b>Hint:</b> Try a general category e.g. <em>history</em> or <em>movies</em>, the title of a podcast, or the name of a channel e. g. <em>RadioLab</em>.
        </span>
      );

    return (
      <DocumentTitle title={getTitle('Search podcasts and channels')}>
        <div>
          <form className="form" onSubmit={this.handleSearch.bind(this)}>
            <Input type="search"
                   ref="query"
                   help={help}
                   onClick={this.handleFocus.bind(this)}
                   placeholder="Find a channel or podcast" />
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

const mapStateToProps = state => {

  const { isLoading } = state.podcasts;
  const { query } = state.search;
  const { isLoggedIn } = state.auth;

  const podcasts = podcastsSelector(state);
  const { channels } = channelsSelector(state);

  return {
    searchQuery: query,
    podcasts: podcasts,
    channels: channels,
    isLoading,
    isLoggedIn
  };
};

export default connect(mapStateToProps)(Search);
