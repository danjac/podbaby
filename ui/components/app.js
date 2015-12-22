import React, { PropTypes } from 'react';
import { Link } from 'react-router';

import 'bootswatch/paper/bootstrap.min.css';

import {
  Nav,
  NavItem,
  Navbar,
  Glyphicon
} from 'react-bootstrap';

import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

import * as actions from '../actions';

const MainNav = props => {
  const { auth } = props;
  const { createHref, isActive } = props.history;
  return (
    <Navbar fixedTop={true}>
      <Navbar.Header>
        <Navbar.Brand>
          <Link to={auth.isLoggedIn ? "/podcasts/" : "/" }><Glyphicon glyph="headphones" /> PodBaby</Link>
        </Navbar.Brand>
      </Navbar.Header>

      {auth.isLoggedIn ?

      <Nav pullLeft={true}>
        <NavItem active={isActive("/podcasts/")} href={createHref("/podcasts/")}><Glyphicon glyph="flash" /> New episodes</NavItem>
        <NavItem active={isActive("/podcasts/list/")} href="#"><Glyphicon glyph="list" /> My podcasts</NavItem>
        <NavItem href="#"><Glyphicon glyph="plus" /> Add new podcast</NavItem>
      </Nav> : ''}

      {auth.isLoggedIn ?

      <Nav pullRight={true}>
        <NavItem href="#"><Glyphicon glyph="cog" /> Settings</NavItem>
        <NavItem href="#" onClick={props.logout}><Glyphicon glyph="log-out" /> Logout</NavItem>
      </Nav> :
      <Nav pullRight={true}>
        <NavItem active={isActive("/login/")}
                 href={createHref("/login/")}><Glyphicon glyph="log-in" /> Login</NavItem>
        <NavItem active={isActive("/signup/")}
                 href={createHref("/signup/")}><Glyphicon glyph="user" /> Signup</NavItem>
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
