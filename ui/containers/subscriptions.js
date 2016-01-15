import _ from 'lodash';
import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import moment from 'moment';
import DocumentTitle from 'react-document-title';

import {
  Input,
  Pagination,
} from 'react-bootstrap';

import * as actions from '../actions';
import { channelsSelector } from '../selectors';
import Icon from '../components/icon';
import Loading from '../components/loading';
import ChannelItem from '../components/channel_item';
import { getTitle } from './utils';

export class Subscriptions extends React.Component {

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

  render() {
    const { page, unfilteredChannels, isLoading } = this.props;

    if (isLoading) {
      return <Loading />;
    }

    if (_.isEmpty(unfilteredChannels) && !isLoading) {
      return (
        <span>You haven't subscribed to any channels yet.
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
      <DocumentTitle title={getTitle('My subscriptions')}>
      <div>
        <Input
          className="form-control"
          type="search"
          ref="filter"
          onClick={this.handleSelect}
          onKeyUp={this.handleFilterChannels}
          placeholder="Find a subscription"
        />
        <Input>
          <a
            className="btn btn-default form-control"
            href={`/podbaby-${moment().format('YYYY-MM-DD')}.opml`}
            download
          ><Icon icon="download" /> Download OPML</a>
        </Input>
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

Subscriptions.propTypes = {
  channels: PropTypes.array.isRequired,
  dispatch: PropTypes.func.isRequired,
  page: PropTypes.object.isRequired,
  isLoading: PropTypes.bool.isRequired,
  unfilteredChannels: PropTypes.array.isRequired,
};

const mapStateToProps = state => {
  return Object.assign({}, channelsSelector(state), { isLoading: state.channels.isLoading });
};

export default connect(mapStateToProps)(Subscriptions);
