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
import { sanitize } from './utils';

const ListItem = props => {
  const {
    podcast,
    createHref,
    isCurrentlyPlaying,
    setCurrentlyPlaying,
    bookmark
  } = props;
  // tbd get audio ref, set played at to last time
  return (
    <div>
      <div className="media">
        <div className="media-body">
          <Grid>
            <Row>
              <Col xs={6} md={6}>
                <h4 className="media-heading">{podcast.title}</h4>
              </Col>
              <Col xs={6} mdPush={3} md={3}>
                <ButtonGroup>
                  <Button onClick={setCurrentlyPlaying}><Glyphicon glyph={ isCurrentlyPlaying ? 'stop': 'play' }  /> </Button>
                  <a className="btn btn-default" href={podcast.enclosureUrl}><Glyphicon glyph="download" /> </a>
                  <Button onClick={bookmark} title={podcast.isBookmarked ? 'Remove bookmark' : 'Add to bookmarks'}>
                    <Glyphicon glyph={podcast.isBookmarked ? 'remove' : 'pushpin'} />
                  </Button>
                </ButtonGroup>
              </Col>
            </Row>
          </Grid>
        </div>
      </div>
      {podcast.description ? <Well dangerouslySetInnerHTML={sanitize(podcast.description)} /> : ''}
    </div>
  );
};

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

  render() {
    const { channel, dispatch, player } = this.props;
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
            <Well>{channel.description}</Well>
          </div>
          </div>
          {channel.podcasts.map(podcast => {
          const isCurrentlyPlaying = player.podcast && podcast.id === player.podcast.id;
          const setCurrentlyPlaying = event => {
            event.preventDefault();
            dispatch(actions.player.setPodcast(isCurrentlyPlaying ? null : podcast));
          };
          const bookmark = (event) => {
            event.preventDefault();
            const { bookmarks } = actions;
            const action = podcast.isBookmarked ? bookmarks.deleteBookmark : bookmarks.addBookmark;
            dispatch(action(podcast.id));
          };
          return <ListItem key={podcast.id}
                           podcast={podcast}
                           bookmark={bookmark}
                           isCurrentlyPlaying={isCurrentlyPlaying}
                           setCurrentlyPlaying={setCurrentlyPlaying}
                           channel={channel} />;
        })}
      </div>
    );
  }
}

Channel.propTypes = {
  channel: PropTypes.object,
  player: PropTypes.object,
  dispatch: PropTypes.func.isRequired
};

const mapStateToProps = state => {
  return {
    channel: state.channel,
    player: state.player
  };
};

export default connect(mapStateToProps)(Channel);
