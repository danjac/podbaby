import React from 'react';
import { Link } from 'react-router';

import {
  Jumbotron,
} from 'react-bootstrap';

const Front = () => {
  const styles = {
    minHeight: 400,
    minWidth: '100%',
    marginTop: 0,
    marginBottom: 8,
    backgroundColor: 'white',
  };

  return (
    <Jumbotron style={styles}>
      <div className="text-center">
        <h1 style={{ fontFamily: 'GoodDog' }}>Podcast lover happiness.</h1>
      <p>
        PodBaby makes it easy to find and listen to your favorite podcasts, and to disover
        new sources of podcasting nirvana.
      </p>
      <p><Link className="btn btn-lg btn-primary" to="/new/">Browse podcasts</Link></p>
      <p><Link className="btn btn-lg btn-success" to="/signup/">Join now</Link></p>
      <p><Link to="/login/">Already a member? Log in here.</Link></p>
      </div>

    </Jumbotron>
  );
};

export default Front;
