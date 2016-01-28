import _ from 'lodash';
import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import DocumentTitle from 'react-document-title';
import { Link } from 'react-router';

import {
  ButtonGroup,
  Button,
  Label,
  Input,
  Grid,
  Row,
  Col,
} from 'react-bootstrap';


import * as actions from '../actions';
import {
  podcastsSelector,
  channelSelector,
  relatedChannelsSelector,
  } from '../selectors';
import PodcastList from '../components/podcasts';
import Image from '../components/image';
import Icon from '../components/icon';
import Loading from '../components/loading';
import { sanitize } from '../components/utils';
import { getTitle } from './utils';

const RelatedChannel = props => {
  const {
    channel,
    handleSubscribe,
    isLoggedIn } = props;

  return (
      <div className="thumbnail">
        <div className="caption text-center">
          <Link to={`/channel/${channel.id}/`}>
            <h5>{channel.title}</h5>
          </Link>
        </div>
        <Image
          src={channel.image}
          errSrc="/static/podcast.png"
          imgProps={{
            alt: channel.title,
            height: 120,
            width: 120,
          }}
        />
        {isLoggedIn ?
          <div className="caption text-center">
            <Button
              title={channel.isSubscribed ? 'Unsubscribe' : 'Subscribe'}
              onClick={handleSubscribe}
            >
            <Icon icon={channel.isSubscribed ? 'unlink' : 'link'} /> {
            channel.isSubscribed ? 'Unsubscribe' : 'Subscribe'
            }
            </Button>
          </div> : ''}
     </div>
  );
};

RelatedChannel.propTypes = {
  channel: PropTypes.object.isRequired,
  handleSubscribe: PropTypes.func.isRequired,
  isLoggedIn: PropTypes.bool.isRequired,
};

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

  handleSelectPage(page) {
    const { channel } = this.props;
    this.actions.channel.getChannel(channel.id, page);
  }

  render() {
    const {
      channel,
      categories,
      isChannelLoading,
      isPodcastsLoading,
      relatedChannels,
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

    const recommendations = relatedChannels.length > 0 ?
    <div className="container">
        <h4 className="text-center">People who subscribed to this feed also subscribed to</h4>
        <Grid>
          <Row>
          {this.props.relatedChannels.map(related => {
            const handleSubscribe = () => this.actions.subscribe.toggleSubscribe(related);
            return (
            <Col key={related.id} xs={12} md={4}>
              <RelatedChannel
                channel={related}
                handleSubscribe={handleSubscribe}
                isLoggedIn={isLoggedIn}
              />
            </Col>
          );
          })}
          </Row>
        </Grid>
    </div> : '';

    return (
      <DocumentTitle title={getTitle(channel.title)}>
        <div>
          <div className="thumbnail">
            <div className="caption text-center">
              <h2>{channel.title}</h2>
            </div>
            <Image
              src={channel.image}
              errSrc="/static/podcast.png"
              imgProps={{
                height: 120,
                width: 120,
                alt: channel.title,
              }}
            />
            {channel.numPodcasts ?
            <div className="caption text-center">
              <h4>
                <Label bsStyle="primary">
                  {channel.numPodcasts} podcast{channel.numPodcasts > 1 ? 's' : ''}
                </Label>
              </h4>
            </div> : ''}
          </div>
          {channel.description ?
          <p
            className="lead"
            style={{ marginTop: 20 }}
            dangerouslySetInnerHTML={sanitize(channel.description)}
          /> : ''}
          <div className="text-center">
            <ButtonGroup>
              {isLoggedIn ?
              <Button
                title={isSubscribed ? 'Unsubscribe' : 'Subscribe'}
                onClick={this.handleSubscribe}
              >
              <Icon icon={isSubscribed ? 'unlink' : 'link'} /> {
              isSubscribed ? 'Unsubscribe' : 'Subscribe'
              }
              </Button> : ''}
              <a
                className="btn btn-default"
                title="Link to RSS Feed"
                target="_blank"
                href={channel.url}
              >
                <Icon icon="rss" /> RSS feed
              </a>
              {website ? (
              <a className="btn btn-default"
                title="Link to home page"
                target="_blank" href={website}
              >
                <Icon icon="globe" /> Website
              </a>
              ) : ''}
            </ButtonGroup>
          </div>
          <div className="text-center" style={{ marginTop: 20 }}>
            <ButtonGroup>
            {categories.map(category => {
              return (
              <Link
                key={category.id}
                className="btn btn-info"
                to={`/categories/${category.id}/`}
              >{category.name}</Link>
              );
            })}
            </ButtonGroup>
        </div>
        <hr />
        <form onSubmit={this.handleSearch}>
          <Input
            type="search"
            ref="query"
            defaultValue={query}
            onClick={this.handleSelectSearch}
            placeholder="Find a podcast from this feed"
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
        {recommendations}
      </div>
      </DocumentTitle>
    );
  }
}

Channel.propTypes = {
  channel: PropTypes.object,
  categories: PropTypes.array,
  relatedChannels: PropTypes.array,
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
  const { query, categories } = state.channel;
  const { page } = state.podcasts;
  const isChannelLoading = state.channel.isLoading;
  const isPodcastsLoading = state.podcasts.isLoading;
  const { isLoggedIn } = state.auth;
  const podcasts = podcastsSelector(state);
  const relatedChannels = relatedChannelsSelector(state);
  const channel = channelSelector(state);

  return {
    podcasts,
    channel,
    categories,
    relatedChannels,
    query,
    page,
    isChannelLoading,
    isPodcastsLoading,
    isLoggedIn,
  };
};

export default connect(mapStateToProps)(Channel);
