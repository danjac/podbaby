import React from 'react';

import {
  Grid,
  Row,
  Col,
  Button,
  ButtonGroup,
  Glyphicon
} from 'react-bootstrap';

const SAMPLE_DATA = [

  {
    image: 'https://gpodder.net/logo/32/341/3419c0f511f571af904efe172acedcf411d07502',
    name: 'Joe Rogan Experience',
    description: 'The podcast of Comedian Joe Rogan.',
    episodes: 745,
    id: 1000
  },
];

const ListItem = props => {
  const { channel } = props;
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
        <h4 className="media-heading"><a href="#">{channel.name}</a></h4>
        <Grid>
          <Row>
            <Col xs={6} md={9}>
              <p>{channel.description}</p>
            </Col>
            <Col xs={6} md={3}>
              <ButtonGroup>
                <Button title="Unsubscribe"><Glyphicon glyph="trash" /></Button>
              </ButtonGroup>
            </Col>
          </Row>
        </Grid>
        <p>
          Episodes: <b>{channel.episodes}</b>
        </p>
      </div>
    </div>
  );
};


export class SubscriptionList extends React.Component {
  render() {
    return (
      <div>
      {SAMPLE_DATA.map(channel => {
        return <ListItem key={channel.id} channel={channel} />;
      })}
      </div>
    );
  }
}

export default SubscriptionList;
