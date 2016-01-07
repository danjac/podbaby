import React from 'react';
import { Link } from 'react-router';

import {
  Nav,
  NavItem,
  Navbar,
  Badge,
  Alert,
  Grid,
  Row,
  Col
} from 'react-bootstrap';

import Icon from './icon';

class NavBar extends React.Component {

  render() {

    const { isLoggedIn, name } = this.props.auth;
    const { createHref, isActive } = this.props.history;

    return (
      <Navbar fixedTop inverse>
        <Navbar.Header>
          <Navbar.Brand>
            <Link style={{ fontFamily: "GoodDog" }} to="/"><Icon icon="headphones" /> PodBaby</Link>
          </Navbar.Brand>
          <Navbar.Toggle />
        </Navbar.Header>

        <Navbar.Collapse>
          <Nav pullLeft>
            <NavItem active={isActive("/new/")}
                     href={createHref("/new/")}><Icon icon="flash" /> New episodes</NavItem>
            <NavItem active={isActive("/search/")}
              href={createHref("/search/")}><Icon icon="search" /> Search</NavItem>
          </Nav>
          {isLoggedIn ?
          <Nav pullLeft>
            <NavItem active={isActive("/member/subscriptions/")}
                     href={createHref("/member/subscriptions/")}><Icon icon="list" /> Subscriptions</NavItem>
            <NavItem active={isActive("/member/bookmarks/")}
                     href={createHref("/member/bookmarks/")}><Icon icon="bookmark" /> Bookmarks</NavItem>
            <NavItem active={isActive("/member/recent/")}
                     href={createHref("/member/recent/")}><Icon icon="clock-o" /> Recently played</NavItem>
            <NavItem onClick={this.props.onOpenAddChannelForm} href="#"><Icon icon="plus" /> Add a channel</NavItem>
          </Nav>
          : ''}

          {isLoggedIn ?
          <Nav pullRight>
            <NavItem active={isActive("/user/")}
                      href={createHref("/user/")}><Icon icon="cog" /> {name}</NavItem>
            <NavItem href="#" onClick={this.props.onLogout}><Icon icon="sign-out" /> Logout</NavItem>
          </Nav> :
          <Nav pullRight>
            <NavItem active={isActive("/login/")}
                     href={createHref("/login/")}><Icon icon="sign-in" /> Login</NavItem>
            <NavItem active={isActive("/signup/")}
                     href={createHref("/signup/")}><Icon icon="sign-in" /> Signup</NavItem>
          </Nav>}
        </Navbar.Collapse>

      </Navbar>
    );
  }
}


export default NavBar;

