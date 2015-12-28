import React, { PropTypes } from 'react';
import { bindActionCreators } from 'react';
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

import { bookmarks, player, subscribe } from '../actions';

import { sanitize } from './utils';

const ListItem = props => {
  const { podcast, createHref, isCurrentlyPlaying, setCurrentlyPlaying, subscribe, deleteBookmark } = props;
  const url = createHref("/podcasts/channel/" + podcast.channelId + "/")
  // tbd get audio ref, set played at to last time
  return (
    <div>
      <div className="media">
        <div className="media-left media-middle">
          <a href={url}>
            <img className="media-object"
                 height={60}
                 width={60}
                 src={podcast.image}
                 alt={podcast.name} />
          </a>
        </div>
        <div className="media-body">
          <h4 className="media-heading"><a href={url}>{podcast.name}</a></h4>
          <Grid>
            <Row>
              <Col xs={6} md={6}>
                <h5>{podcast.title}</h5>
              </Col>
              <Col xs={6} mdPush={3} md={3}>
                <ButtonGroup>
                  <Button onClick={setCurrentlyPlaying}><Glyphicon glyph={ isCurrentlyPlaying ? 'stop': 'play' }  /></Button>
                    <a title="Download this podcast" className="btn btn-default" href={podcast.enclosureUrl}><Glyphicon glyph="download" /></a>
                  <Button onClick={deleteBookmark} title="Remove this bookmark"><Glyphicon glyph="remove" /></Button>
                  <Button title={podcast.isSubscribed ? "Unsubscribe" : "Subscribe"} onClick={subscribe}>
                    <Glyphicon glyph={podcast.isSubscribed ? "trash" : "ok"} />
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


export class Bookmarks extends React.Component {

  componentDidMount() {
    const { dispatch } = this.props;
    dispatch(bookmarks.getBookmarks());
  }

  render() {
    const { dispatch } = this.props;
    const podcasts = this.props.bookmarks;
    const { createHref } = this.props.history;
    if (podcasts.length === 0) {
      return <div>You do not have any bookmarked podcasts yet.</div>;
    }
    return (
      <div>
        {podcasts.map(podcast => {
          const isCurrentlyPlaying = this.props.player.podcast && podcast.id === this.props.player.podcast.id;
          const setCurrentlyPlaying = (event) => {
            event.preventDefault();
            dispatch(player.setPodcast(isCurrentlyPlaying ? null : podcast));
          };
          const doSubscribe = event => {
            event.preventDefault();
            const action = podcast.isSubscribed ? subscribe.unsubscribe : subscribe.subscribe;
            dispatch(action(podcast.channelId, podcast.name));
          };
          const deleteBookmark = event => {
            event.preventDefault();
            dispatch(bookmarks.deleteBookmark(podcast.id));
          };
          return <ListItem key={podcast.id}
                           podcast={podcast}
                           subscribe={doSubscribe}
                           deleteBookmark={deleteBookmark}
                           isCurrentlyPlaying={isCurrentlyPlaying}
                           setCurrentlyPlaying={setCurrentlyPlaying}
                           createHref={createHref} />;
        })}
      </div>
    );
  }
}

Bookmarks.propTypes = {
  bookmarks: PropTypes.array.isRequired,
  currentlyPlaying: PropTypes.number,
  dispatch: PropTypes.func.isRequired
};

const mapStateToProps = state => {
  return {
    bookmarks: state.bookmarks,
    player: state.player
  };
};

export default connect(mapStateToProps)(Bookmarks);
