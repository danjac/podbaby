import React, { PropTypes } from 'react';

import {
  Input,
  Button
} from 'react-bootstrap';

export class Signup extends React.Component {
  render() {
    return (
      <div>
        <form className="form-horizontal">
            <Input required
              type="text"
              ref="username"
              pattern="[a-zA-Z0-9]{3}"
              placeholder="Username" />
            <Input required
              type="email"
              ref="email"
              placeholder="Email address" />
            <Input required
              type="password"
              ref="password"
              placeholder="Password" />
            <Button
              bsStyle="primary"
              className="form-control"
              type="submit">Signup</Button>
        </form>
      </div>

    );
  }
};

export default Signup;
