import React, { PropTypes } from 'react';
import { connect } from 'react-redux';

import {
  Grid,
  Row,
  Col,
  ButtonGroup,
  Button
} from 'react-bootstrap';

import * as actions from '../actions';
import PodcastList from './podcasts';
import Icon from './icon';
import Loading from './loading';
import { sanitize, formatPubDate } from './utils';

export class Channel extends React.Component {

  handleSubscribe(event) {
    event.preventDefault();
    const { channel, dispatch } = this.props;
    const action = channel.isSubscribed ? actions.subscribe.unsubscribe : actions.subscribe.subscribe;
    dispatch(action(channel.id));
  }

  handleSelectPage(event, selectedEvent) {
    event.preventDefault();
    const { dispatch } = this.props;
    const page = selectedEvent.eventKey;
    dispatch(actions.channel.getChannel(this.props.params.id, page));
  }

  render() {
    const { channel, isLoading } = this.props;

    if (isLoading) {
      return <Loading />;
    }

    if (!channel) {
      return <div>Sorry, could not find this channel.</div>;
    }
    const isSubscribed = channel.isSubscribed;

    const website = channel.website.Valid ? channel.website.String : "";

    return (
      <div>
        <div className="media">
          <div className="media-left">
            <a href="#">
              <img className="media-object"
                   height={60}
                   width={60}
                   src={channel.image}
                   alt={channel.title} />
            </a>
          </div>
          <div className="media-body">
                  <h2 className="media-heading">{channel.title}</h2>
          </div>
        </div>
        {channel.description ? <p className="lead" style={{ marginTop: 20 }} dangerouslySetInnerHTML={sanitize(channel.description)} /> : ''}
        <ButtonGroup>
          <Button title={channel.isSubscribed ? 'Unsubscribe': 'Subscribe'}
                  onClick={this.handleSubscribe.bind(this)}>
            <Icon icon={channel.isSubscribed ? 'unlink': 'link'} /> {channel.isSubscribed ? 'Unsubscribe' : 'Subscribe'}</Button>
          <a className="btn btn-default" title="Link to RSS Feed" target="_blank" href={channel.url}>
            <Icon icon="rss" /> Link to RSS feed
          </a>
          {website ? (
          <a className="btn btn-default" title="Link to home page" target="_blank" href={website}>
            <Icon icon="globe" /> Link to website
          </a>
          ) : ''}
        </ButtonGroup>
        <hr />
        <PodcastList showChannel={false}
                     onSelectPage={this.handleSelectPage.bind(this)}
                     actions={actions} {...this.props} />
      </div>
    );
  }
}

Channel.propTypes = {
  channel: PropTypes.object,
  podcasts: PropTypes.array,
  page: PropTypes.object,
  player: PropTypes.object,
  dispatch: PropTypes.func.isRequired
};

const mapStateToProps = state => {

  const { channel } = state.channel;
  const { podcasts, page, showDetail } = state.podcasts;
  const isLoading = state.channel.isLoading || state.podcasts.isLoading;

  return {
    player: state.player,
    podcasts: podcasts || [],
    channel,
    showDetail,
    isLoading,
    page
  };
};

export default connect(mapStateToProps)(Channel);
