import React, { PropTypes } from 'react';
import { connect } from 'react-redux';

import {
  Grid,
  Row,
  Col,
  Glyphicon,
  ButtonGroup,
  Button
} from 'react-bootstrap';

import * as actions from '../actions';
import PodcastList from './podcasts';
import Loading from './loading';
import { sanitize, formatPubDate } from './utils';

export class Channel extends React.Component {

  componentDidMount(){
      this.props.dispatch(actions.channel.getChannel(this.props.params.id));
  }

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
            <Grid>
              <Row>
                <Col xs={6} md={9}>
                  <h2 className="media-heading">{channel.title}</h2>
                </Col>
                <Col xs={6} md={3}>
                  <ButtonGroup>
                    <Button title={channel.isSubscribed ? 'Unsubscribe': 'Subscribe'}
                            onClick={this.handleSubscribe.bind(this)}>
                      <Glyphicon glyph={channel.isSubscribed ? 'trash': 'ok'} /> {channel.isSubscribed ? 'Unsubscribe' : 'Subscribe'}</Button>
                  </ButtonGroup>
                </Col>
              </Row>
            </Grid>
          </div>
          </div>
          {channel.description ? <p className="lead" style={{ marginTop: 20 }} dangerouslySetInnerHTML={sanitize(channel.description)} /> : ''}
          <p>
            <a href={channel.url}>RSS Feed</a> {channel.website? <a target="_blank" href={channel.website}>Website</a> : ''}
          </p>
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
