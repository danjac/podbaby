import _ from 'lodash';
import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import DocumentTitle from 'react-document-title';

import {
  ButtonGroup,
  Button,
  Input,
} from 'react-bootstrap';


import * as actions from '../actions';
import { podcastsSelector, channelSelector } from '../selectors';
import PodcastList from '../components/podcasts';
import Image from '../components/image';
import Icon from '../components/icon';
import Loading from '../components/loading';
import { sanitize } from '../components/utils';
import { getTitle } from './utils';

export class Channel extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = {
      channel: bindActionCreators(actions.channel, dispatch),
      subscribe: bindActionCreators(actions.subscribe, dispatch),
    };
    this.handleSelectPage = this.handleSelectPage.bind(this);
    this.handleSearch = this.handleSearch.bind(this);
    this.handleClearSearch = this.handleClearSearch.bind(this);
    this.handleSelectSearch = this.handleSelectSearch.bind(this);
    this.handleSubscribe = this.handleSubscribe.bind(this);
  }

  handleSearch(event) {
    event.preventDefault();
    const { channel } = this.props;
    const value = this.refs.query.getValue();
    const query = _.trim(value);
    if (query) {
      this.actions.channel.searchChannel(query, channel.id);
    } else {
      this.actions.channel.getChannel(channel.id);
    }
  }

  handleSelectSearch(event) {
    event.preventDefault();
    this.refs.query.getInputDOMNode().select();
  }

  handleClearSearch(event) {
    event.preventDefault();
    const { channel } = this.props;
    this.actions.channel.getChannel(channel.id);
    this.refs.query.getInputDOMNode().value = '';
  }

  handleSubscribe(event) {
    event.preventDefault();
    const { channel } = this.props;
    this.actions.subscribe.toggleSubscribe(channel);
  }

  handleSelectPage(event, selectedEvent) {
    event.preventDefault();
    const { channel } = this.props;
    const page = selectedEvent.eventKey;
    this.actions.channel.getChannel(channel.id, page);
  }

  render() {
    const {
      channel,
      isChannelLoading,
      isPodcastsLoading,
      query,
      isLoggedIn } = this.props;

    if (isChannelLoading) {
      return <Loading />;
    }

    if (!channel) {
      return <div>Sorry, could not find this channel.</div>;
    }

    const website = channel.website.Valid ? channel.website.String : '';
    const { isSubscribed } = channel;

    return (
      <DocumentTitle title={getTitle(channel.title)}>
        <div>
        <div className="media">
          <div className="media-left">
            <a href="#">
              <Image
                className="media-object"
                src={channel.image}
                errSrc="/static/podcast.png"
                imgProps={{
                  height: 60,
                  width: 60,
                  alt: channel.title,
                }}
              />
            </a>
          </div>
          <div className="media-body">
                  <h2 className="media-heading">{channel.title}</h2>
          </div>
        </div>
        {channel.description ?
        <p
          className="lead"
          style={{ marginTop: 20 }}
          dangerouslySetInnerHTML={sanitize(channel.description)}
        /> : ''}
        <ButtonGroup>
          {isLoggedIn ?
          <Button
            title={isSubscribed ? 'Unsubscribe' : 'Subscribe'}
            onClick={this.handleSubscribe}
          >
            <Icon icon={isSubscribed ? 'unlink' : 'link'} />
            {isSubscribed ? 'Unsubscribe' : 'Subscribe'}
          </Button> : ''}
          <a
            className="btn btn-default"
            title="Link to RSS Feed"
            target="_blank"
            href={channel.url}
          >
            <Icon icon="rss" /> Link to RSS feed
          </a>
          {website ? (
          <a className="btn btn-default" title="Link to home page" target="_blank" href={website}>
            <Icon icon="globe" /> Link to website
          </a>
          ) : ''}
        </ButtonGroup>
        <hr />
        <form onSubmit={this.handleSearch}>
          <Input
            type="search"
            ref="query"
            defaultValue={query}
            onClick={this.handleSelectSearch}
            placeholder="Find a podcast in this channel"
          />
          <Input>
            <Button
              bsStyle="primary"
              type="submit"
              className="form-control"
            >
              <Icon icon="search" /> Search
            </Button>
          </Input>
          {query ?
          <Input>
            <Button
              bsStyle="default"
              onClick={this.handleClearSearch}
              className="form-control"
            >
              <Icon icon="refresh" /> Show all podcasts
            </Button>
          </Input> : ''}
        </form>
        {isPodcastsLoading && !query ? <Loading /> :
        <PodcastList
          showChannel={false}
          isLoggedIn={isLoggedIn}
          isLoading={isPodcastsLoading}
          onSelectPage={this.handleSelectPage}
          actions={actions}
          {...this.props}
        /> }
      </div>
      </DocumentTitle>
    );
  }
}

Channel.propTypes = {
  channel: PropTypes.object,
  podcasts: PropTypes.array,
  page: PropTypes.object,
  player: PropTypes.object,
  dispatch: PropTypes.func.isRequired,
  isChannelLoading: PropTypes.bool.isRequired,
  isPodcastsLoading: PropTypes.bool.isRequired,
  isLoggedIn: PropTypes.bool.isRequired,
  query: PropTypes.string,
};

const mapStateToProps = state => {
  const { query } = state.channel;
  const { page } = state.podcasts;
  const isChannelLoading = state.channel.isLoading;
  const isPodcastsLoading = state.podcasts.isLoading;
  const { isLoggedIn } = state.auth;
  const podcasts = podcastsSelector(state);
  const channel = channelSelector(state);

  return {
    podcasts,
    channel,
    query,
    page,
    isChannelLoading,
    isPodcastsLoading,
    isLoggedIn,
  };
};

export default connect(mapStateToProps)(Channel);
