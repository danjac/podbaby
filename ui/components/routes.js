import React, { PropTypes } from 'react';
import { pushPath, replacePath } from 'redux-simple-router';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Router, Route, IndexRoute } from 'react-router';

import App from './app';
import Dashboard from './dashboard';
import Front from './front';
import Login from './login';
import Signup from './signup';
import PageNotFound from './not_found';

export class Routes extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindActionCreators({ pushPath }, dispatch);
  }

  render() {

    const loginRequired = (nextState, replaceState) => {
      // we'll make a custom action later
//    // if (!this.props.auth.isLoggedIn )....
      replaceState(null, "/login/");
    };

    return (
      <Router history={this.props.history}>
        <Route path="/" component={App}>
          <Route path="/secure/" onEnter={loginRequired}>
            <IndexRoute component={Dashboard} />
          </Route>
          <IndexRoute component={Front} />
          <Route path="/login/" component={Login} />
          <Route path="/signup/" component={Signup} />
          <Route path="*" component={PageNotFound} />
        </Route>
      </Router>
    );
  }
}

Routes.propTypes = {
  dispatch: PropTypes.func.isRequired
};

export default connect()(Routes);
