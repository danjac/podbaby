import React from 'react';
import { Glyphicon } from 'react-bootstrap';

export default function(props) {

  return (
      <div className="text-center" style={{ marginTop: 50 }}>
        <h1><Glyphicon glyph="refresh" /> Loading...</h1>
      </div>
  );

}
