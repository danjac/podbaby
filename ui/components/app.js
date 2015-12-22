import React, { PropTypes } from 'react';
import { Link } from 'react-router';

import 'bootswatch/paper/bootstrap.min.css';

import {
  Nav,
  NavItem,
  Navbar
} from 'react-bootstrap';

import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

import * as actions from '../actions';

const MainNav = props => {
  const { createHref, isActive } = props.history;
  return (
    <Navbar fixedTop={true} inverse={true}>
      <Navbar.Header>
        <Navbar.Brand>
          <Link to="/">Podbaby</Link>
        </Navbar.Brand>
      </Navbar.Header>
      <Nav pullLeft={true}>
        <NavItem active={isActive("/secure/")}
                 href={createHref("/secure/")}>Dashboard</NavItem>
      </Nav>
      {props.auth.isLoggedIn ?

      <Nav pullRight={true}>
        <NavItem href="#" onClick={props.logout}>Logout</NavItem>
      </Nav> :
      <Nav pullRight={true}>
        <NavItem active={isActive("/login/")}
                 href={createHref("/login/")}>Login</NavItem>
        <NavItem active={isActive("/signup/")}
                 href={createHref("/signup/")}>Signup</NavItem>
      </Nav>}

    </Navbar>
  );
};

export class App extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindActionCreators(actions.auth, dispatch);
  }

  logout(event) {
    event.preventDefault();
    this.actions.logout();
  }

  render() {

    return (
      <div>
        <MainNav logout={this.logout.bind(this)} {...this.props} />
        <div className="container" style={ { marginTop: "80px"}  }>
          {this.props.children}
        </div>
      </div>
    );
  }
}

App.propTypes = {
  routing: PropTypes.object.isRequired,
  dispatch: PropTypes.func.isRequired,
  auth: PropTypes.object
};


const mapStateToProps = state => {
  const { routing, auth } = state;
  return {
    routing,
    auth
  };
};

export default connect(mapStateToProps)(App);
