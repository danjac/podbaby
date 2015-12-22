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

const SAMPLE_DATA = [

  {
    image: 'https://gpodder.net/logo/32/341/3419c0f511f571af904efe172acedcf411d07502',
    name: 'Joe Rogan Experience',
    title: 'Molly Crabapple',
    summary: '#738. Molly Crabapple is an American artist and journalist, known for her work for The New York Times, VICE, the Wall Street Journal, the Royal Society of Arts, Red Bull, Marvel Comics, DC Comics and CNN. Her new book "Drawing Blood" is available now on Amazon',
    id: 1000
  },
  {
    image: 'https://gpodder.net/logo/32/f3f/f3f419da4a5e90d5e7eb1fdfd032dd9c327d2494',
    name: 'Radiolab from NYC',
    title: 'Nazi Summer Camp',
    summary: 'The incredible, little-known story of the Nazi prisoners of war kept on American soil during World War II.',
    id: 1001
  },
  {
    image: 'https://gpodder.net/logo/32/124/1246edc674de518c36549def4493b47eb43d4591',
    name: 'JavaScript Jabber',
    title: ' 146 JSJ React with Christopher Chedeau and Jordan Walker',
    summary: 'The panelists talk to Christopher Chedeau and Jordan Walke about React.js Conf and React Native.',
    id: 1002
  }


];

const ListItem = props => {
  const { episode } = props;
  return (
    <div className="media">
      <div className="media-left">
        <a href="#">
          <img className="media-object"
               src={episode.image}
               alt={episode.title} />
        </a>
      </div>
      <div className="media-body">
        <h4 className="media-heading"><a href="#">{episode.name}</a></h4>
        <Grid>
          <Row>
            <Col xs={6} md={9}>
              <h5>{episode.title}</h5>
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
        <Well>{episode.summary}</Well>
      </div>
    </div>
  );
};

export class EpisodeList extends React.Component {
  render() {
    return (
      <div>
        {SAMPLE_DATA.map(episode => {
          return <ListItem key={episode.id} episode={episode} />;
        })}
      </div>
    );
  }
}

export default EpisodeList;
