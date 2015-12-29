import React, { PropTypes } from 'react';

import {
  Grid,
  Row,
  Col,
  Glyphicon,
  ButtonGroup,
  Button,
  Well
} from 'react-bootstrap';

export const Podcast = props => {

  const {
    podcast,
    showChannel,
    channelUrl,
    isPlaying,
    isShowingDetail,
    togglePlayer,
    toggleSubscribe,
    toggleDetail,
    toggleBookmark } = props;

  const header = showChannel ? <h3><a href={channelUrl}>{podcast.name}</a></h3> : <h3>{podcast.title}</h3>;

  return (
    <Panel header={header}>
      <div className="media">
        {showChannel ?
        (<div className="media-left media-middle">
          <a href={channelUrl}>
            <img className="media-object"
                 height={60}
                 width={60}
                 src={podcast.image}
                 alt={podcast.name} />
          </a>
          </div> ) : ''}
        <div className="media-body">
          <Grid>
            <Row>
              <Col xs={6} md={6}>
                {showChannel ? <h4>{podcast.title}</h4> : ''}
                <b>{formatPubDate(podcast.pubDate)}</b>
              </Col>
              <Col xs={6} mdPush={2} md={3}>
                <ButtonGroup>
                  <Button onClick={togglePlayer}><Glyphicon glyph={ isPlaying ? 'stop': 'play' }  /></Button>
                  <a title="Download this podcast" className="btn btn-default" href={podcast.enclosureUrl}><Glyphicon glyph="download" /></a>
                  <Button onClick={toggleBookmark} title={podcast.isBookmarked ? 'Remove bookmark' : 'Add to bookmarks'}>
                    <Glyphicon glyph={podcast.isBookmarked ? 'remove' : 'bookmark'} />
                  </Button>
                  <Button title="Unsubscribe from this channel" onClick={unsubscribe}><Glyphicon glyph="trash" /></Button>
                </ButtonGroup>
              </Col>
            </Row>
          </Grid>
        </div>
      </div>
      {podcast.description ?
      (<div style={{paddingTop: "30px"}}>
        <Button className="form-control" onClick={toggleDetail}>
        {isShowingDetail ? 'Show less' : 'Show more'} <Glyphicon glyph={isShowingDetail ? 'chevron-up' : 'chevron-down'} />
        </Button>
      </div>) : ''}
      {podcast.description && isShowingDetail  ? <Well dangerouslySetInnerHTML={sanitize(podcast.description)} /> : ''}
  </Panel>
  );
};
