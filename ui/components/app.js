import React, { PropTypes } from 'react';
import { Link } from 'react-router';

import 'bootswatch/paper/bootstrap.min.css';

import {
  Nav,
  NavItem,
  Navbar
} from 'react-bootstrap';

import { connect } from 'react-redux';

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
      <Nav pullRight={true}>
        <NavItem active={isActive("/login/")}
                 href={createHref("/login/")}>Login</NavItem>
        <NavItem active={isActive("/signup/")}
                 href={createHref("/signup/")}>Signup</NavItem>
      </Nav>
    </Navbar>
  );
};

export class App extends React.Component {
  render() {
    return (
      <div>
        <MainNav {...this.props} />
        <div className="container" style={ { marginTop: "80px"}  }>
          {this.props.children}
        </div>
      </div>
    );
  }
}

App.propTypes = {
  routing: PropTypes.object.isRequired
};


const mapStateToProps = props => {
  return {
    routing: props.routing
  };
};

export default connect(mapStateToProps)(App);
