import React, { PropTypes } from 'react';
import { connect } from 'react-redux';

import {
  Grid,
  Row,
  Col,
  Glyphicon,
  ButtonGroup,
  Button,
  Well
} from 'react-bootstrap';

import * as actions from '../actions';
import PodcastList from './podcasts';
import { sanitize, formatPubDate } from './utils';

export class Channel extends React.Component {

  componentDidMount(){
      this.props.dispatch(actions.channel.getChannel(this.props.params.id));
  }

  handleSubscribe(event) {
    event.preventDefault();
    const { channel, dispatch } = this.props;
    const action = channel.isSubscribed ? actions.subscribe.unsubscribe : actions.subscribe.subscribe;
    dispatch(action(channel.id, channel.title));
  }

  handleSelectPage(event, selectedEvent) {
    event.preventDefault();
    const { dispatch } = this.props;
    const page = selectedEvent.eventKey;
    dispatch(actions.channel.getChannel(this.props.params.id, page));
  }

  render() {
    const { channel, isLoading } = this.props;
    if (!channel) {
      return <div></div>;
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
                    <Button title="Unsubscribe" onClick={this.handleSubscribe.bind(this)}><Glyphicon glyph="trash" /> Unsubscribe</Button>
                  </ButtonGroup>
                </Col>
              </Row>
            </Grid>
            {channel.description ? <Well dangerouslySetInnerHTML={sanitize(channel.description)} /> : ''}
          </div>
          </div>
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
  const { podcasts, page, showDetail, isLoading } = state.podcasts;
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
