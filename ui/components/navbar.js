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

    const { auth } = this.props;
    const { createHref, isActive } = this.props.history;

    return (
      <Navbar fixedTop inverse>
        <Navbar.Header>
          <Navbar.Brand>
            <Link style={{ fontFamily: "GoodDog" }} to="/podcasts/new/"><Icon icon="headphones" /> PodBaby</Link>
          </Navbar.Brand>
          <Navbar.Toggle />
        </Navbar.Header>

        <Navbar.Collapse>
          <Nav pullLeft>
            <NavItem active={isActive("/podcasts/search/")}
              href={createHref("/podcasts/search/")}><Icon icon="search" /> Search</NavItem>
            <NavItem active={isActive("/podcasts/new/")}
                     href={createHref("/podcasts/new/")}><Icon icon="flash" /> New episodes</NavItem>
            <NavItem active={isActive("/podcasts/subscriptions/")}
                     href={createHref("/podcasts/subscriptions/")}><Icon icon="list" /> Subscriptions</NavItem>
            <NavItem active={isActive("/podcasts/bookmarks/")}
                     href={createHref("/podcasts/bookmarks/")}><Icon icon="bookmark" /> Bookmarks</NavItem>
            <NavItem active={isActive("/podcasts/recent/")}
                     href={createHref("/podcasts/recent/")}><Icon icon="clock-o" /> Recently played</NavItem>
            <NavItem onClick={this.props.onOpenAddChannelForm} href="#"><Icon icon="plus" /> Add a channel</NavItem>
          </Nav>

          <Nav pullRight>
            <NavItem active={isActive("/user/")}
                      href={createHref("/user/")}><Icon icon="cog" /> {auth.name}</NavItem>
            <NavItem href="#" onClick={this.props.onLogout}><Icon icon="sign-out" /> Logout</NavItem>
          </Nav>
        </Navbar.Collapse>

      </Navbar>
    );
  }
}


export default NavBar;

