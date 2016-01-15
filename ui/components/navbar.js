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

  constructor(props) {
    super(props);
    this.state = { expanded: false };
  }

  handleToggle() {
    this.setState({ expanded: !this.state.expanded });
  }

  handleSelected() {
    this.setState({ expanded: false });
  }

  handleOpenAddChannelForm(event) {
    this.props.onOpenAddChannelForm(event);
    this.handleSelected();
  }

  render() {

    const { isLoggedIn, name } = this.props.auth;
    const { createHref, isActive } = this.props.history;

    const handleSelected = this.handleSelected.bind(this);

    return (
      <Navbar fixedTop
              expanded={this.state.expanded}
              onToggle={this.handleToggle.bind(this)}>
        <Navbar.Header>
          <Navbar.Brand>
            <Link style={{ fontFamily: "GoodDog" }} to="/"><Icon icon="headphones" /> PodBaby</Link>
          </Navbar.Brand>
          <Navbar.Toggle />
        </Navbar.Header>

        <Navbar.Collapse>
          <Nav pullLeft>
            <NavItem active={isActive("/new/")}
                     href={createHref("/new/")}
                     onClick={handleSelected}><Icon icon="flash" /> New episodes</NavItem>
            <NavItem active={isActive("/search/")}
                     href={createHref("/search/")}
                     onClick={handleSelected}><Icon icon="search" /> Search</NavItem>
          </Nav>
          {isLoggedIn ?
          <Nav pullLeft>
            <NavItem active={isActive("/member/subscriptions/")}
                     href={createHref("/member/subscriptions/")}
                     onClick={handleSelected}><Icon icon="list" /> Subscriptions</NavItem>
            <NavItem active={isActive("/member/bookmarks/")}
                     href={createHref("/member/bookmarks/")}
                     onClick={handleSelected}><Icon icon="bookmark" /> Bookmarks</NavItem>
            <NavItem active={isActive("/member/recent/")}
                     href={createHref("/member/recent/")}
                     onClick={handleSelected}><Icon icon="clock-o" /> Recently played</NavItem>
            <NavItem onClick={this.handleOpenAddChannelForm.bind(this)}
                     href="#"><Icon icon="rss" /> Add RSS feed</NavItem>
          </Nav>
          : ''}

          {isLoggedIn ?
          <Nav pullRight>
            <NavItem active={isActive("/user/")}
                     href={createHref("/user/")}
                     onClick={handleSelected}><Icon icon="cog" /> {name}</NavItem>
            <NavItem href="#" onClick={this.props.onLogout}><Icon icon="sign-out" /> Logout</NavItem>
          </Nav> :
          <Nav pullRight>
            <NavItem active={isActive("/login/")}
                     href={createHref("/login/")}
                     onClick={handleSelected}><Icon icon="sign-in" /> Login</NavItem>
            <NavItem active={isActive("/signup/")}
                     href={createHref("/signup/")}
                     onClick={handleSelected}><Icon icon="sign-in" /> Signup</NavItem>
          </Nav>}
        </Navbar.Collapse>

      </Navbar>
    );
  }
}


export default NavBar;

