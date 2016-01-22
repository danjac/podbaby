import _ from 'lodash';
import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import DocumentTitle from 'react-document-title';

import {
  Input,
  Pagination,
  ButtonGroup,
  Breadcrumb,
  BreadcrumbItem,
} from 'react-bootstrap';

import * as actions from '../actions';
import { channelsSelector, categorySelector } from '../selectors';
import { isMobile } from '../components/utils';
import Loading from '../components/loading';
import ChannelItem from '../components/channel_item';
import { getTitle } from './utils';

export class Category extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindActionCreators(actions.channels, dispatch);
    this.handleFilterChannels = this.handleFilterChannels.bind(this);
    this.handleSelectPage = this.handleSelectPage.bind(this);
    this.handleSelect = this.handleSelect.bind(this);
  }

  handleFilterChannels() {
    const value = _.trim(this.refs.filter.getValue());
    this.actions.filterChannels(value);
  }

  handleSelectPage(event, selectedEvent) {
    event.preventDefault();
    const page = selectedEvent.eventKey;
    this.actions.selectPage(page);
  }

  handleSelect() {
    this.refs.filter.getInputDOMNode().select();
  }

  renderBreadcrumbs() {
    const { createHref } = this.props.history;
    const { category } = this.props;

    const items = [<BreadcrumbItem key="all" href={createHref('/browse/')}>Browse</BreadcrumbItem>];
    let parent = category.parent;
    while (parent) {
      items.push(
      <BreadcrumbItem key={parent.id} href={createHref(`/categories/${parent.id}/`)}>
      {parent.name}
      </BreadcrumbItem>);
      parent = parent.parent;
    }
    items.push(<BreadcrumbItem key="active" active>{category.name}</BreadcrumbItem>);
    return <Breadcrumb>{items}</Breadcrumb>;
  }

  renderChildren() {
    const { category } = this.props;
    return (
      <ButtonGroup vertical={isMobile()}>
        {category.children.map(child => {
          return (
          <Link
            key={child.id}
            className="btn btn-info"
            to={`/categories/${child.id}/`}
          >{child.name}</Link>
          );
        })}
      </ButtonGroup>
    );
  }

  render() {
    const { page, unfilteredChannels, isLoading, category } = this.props;

    if (isLoading) {
      return <Loading />;
    }

    if (!category || _.isEmpty(unfilteredChannels)) {
      return (
        <span>There are no feeds for this category.
          Discover new channels and podcasts <Link to="/search/">here</Link>.</span>);
    }

    const pagination = (
      page && page.numPages > 1 ?
      <Pagination
        onSelect={this.handleSelectPage}
        first
        last
        prev
        next
        maxButtons={6}
        items={page.numPages}
        activePage={page.page}
      /> : '');

    return (
      <DocumentTitle title={getTitle(`Category: ${category.name}`)}>
      <div>
        {this.renderBreadcrumbs()}
        {this.renderChildren()}
        <Input
          className="form-control"
          type="search"
          ref="filter"
          onClick={this.handleSelect}
          onKeyUp={this.handleFilterChannels}
          placeholder="Find a feed in this category"
        />
        {pagination}
      {this.props.channels.map(channel => {
        const toggleSubscribe = () => {
          this.props.dispatch(actions.subscribe.toggleSubscribe(channel));
        };
        return (
          <ChannelItem
            key={channel.id}
            channel={channel}
            isLoggedIn
            subscribe={toggleSubscribe}
          />
        );
      })}
      </div>
    </DocumentTitle>
    );
  }
}

Category.propTypes = {
  history: PropTypes.object.isRequired,
  channels: PropTypes.array.isRequired,
  category: PropTypes.object.isRequired,
  dispatch: PropTypes.func.isRequired,
  page: PropTypes.object.isRequired,
  isLoading: PropTypes.bool.isRequired,
  unfilteredChannels: PropTypes.array.isRequired,
};

const mapStateToProps = state => {
  return Object.assign({},
    channelsSelector(state), {
      isLoading: state.channels.isLoading,
      category: categorySelector(state),
    });
};

export default connect(mapStateToProps)(Category);
