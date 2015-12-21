import React, { PropTypes } from 'react';

import {
  Input,
  Button
} from 'react-bootstrap';

export class Login extends React.Component {
  render() {
    return (
      <div>
        <form className="form-horizontal">
            <Input required
              type="text"
              ref="identifier"
              placeholder="Email or username" />
            <Input required
              type="password"
              ref="password"
              placeholder="Password" />
            <Button
              bsStyle="primary"
              className="form-control"
              type="submit">Login</Button>
        </form>
      </div>

    );
  }
};

export default Login;
