import React from 'react';
import Icon from './icon';

export default function(props) {

  return (
      <div className="text-center" style={{ marginTop: 50 }}>
        <h1 style={{ fontFamily: "GoodDog"}}><Icon icon="spinner" spin /> loading...</h1>
      </div>
  );

}
