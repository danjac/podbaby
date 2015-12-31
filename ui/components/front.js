import React from 'react';

import {
  Jumbotron
} from 'react-bootstrap';

const Front = props => {
  const { createHref } = props.history;

  const styles = {
    minHeight: 400,
    minWidth: "100%",
    marginTop: 0,
    marginBottom: 8
  };

  return (
    <Jumbotron style={styles}>
      <div className="text-center">
      <h1>Podcast lover happiness.</h1>
      <p>
        PodBaby makes it easy to find and listen to your favorite podcasts, and to disover new sources of podcasting nirvana.
      </p>
      <p><a className="btn btn-lg btn-success" href={createHref('/signup/')}>Join now</a></p>
      <p><a href={createHref('/login/')}>Already a member? Log in here.</a></p>
      </div>

    </Jumbotron>
  );
};

export default Front;
