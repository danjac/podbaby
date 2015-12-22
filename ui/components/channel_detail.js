import React from 'react';

import {
  Grid,
  Row,
  Col,
  Glyphicon,
  ButtonGroup,
  Button,
  Well
} from 'react-bootstrap';

const SAMPLE_DATA = {

  channel: {
    name: 'Joe Rogan Experience',
    image: 'https://gpodder.net/logo/32/341/3419c0f511f571af904efe172acedcf411d07502'
  },

  podcasts: [
    {
      title: 'Molly Crabapple',
      summary: 'Molly Crabapple is an American artist and journalist, known for her work for The New York Times, VICE, the Wall Street Journal, the Royal Society of Arts, Red Bull, Marvel Comics, DC Comics and CNN. Her new book "Drawing Blood" is available now on Amazon',
      id: 1000
    },
    {
      title: 'Brian Stann',
      summary: 'Brian Stann is a retired mixed martial artist and U.S. Marine, who competed as a Middleweight in the UFC. He is a former WEC Light Heavyweight Champion and is currently an analyst & commentator for the UFC.',
      id: 1001
    },
  ]
}

export class ChannelDetail extends React.Component {
  render() {
    const { channel, podcasts } = SAMPLE_DATA;
    return (
      <div>
        <h2>{channel.name}</h2>
        {podcasts.map(podcast => {
        return <ListItem key={podcast.id}
                         podcast={podcast}
                         channel={channel} />;
        })}
      </div>
    );
  }
}

const ListItem = props => {
  const { podcast, channel } = props;
  return (
    <div className="media">
      <div className="media-left">
        <a href="#">
          <img className="media-object"
               src={channel.image}
               alt={channel.name} />
        </a>
      </div>
      <div className="media-body">
        <h4 className="media-heading">{podcast.title}</h4>
        <Grid>
          <Row>
            <Col xs={6} md={9}>
              <p>{podcast.description}</p>
            </Col>
            <Col xs={6} md={3}>
              <ButtonGroup>
                <Button><Glyphicon glyph="play" /></Button>
                <Button><Glyphicon glyph="download" /></Button>
                <Button><Glyphicon glyph="pushpin" /></Button>
                <Button><Glyphicon glyph="ok" /></Button>
              </ButtonGroup>
            </Col>
          </Row>
        </Grid>
      </div>
    </div>
  );
};

export default ChannelDetail;
