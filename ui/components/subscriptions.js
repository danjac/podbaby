import _ from 'lodash';
import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import moment from 'moment';

import {
  Panel,
  Input
} from 'react-bootstrap';

import * as  actions from '../actions';
import { channelsSelector } from '../selectors';
import Icon from './icon';
import Loading from './loading';
import ChannelItem from './channel_item';

export class Subscriptions extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindActionCreators(actions.channels, dispatch);
  }

  handleFilterChannels() {
    const value = _.trim(this.refs.filter.getValue());
    this.actions.filterChannels(value);
  }

  handleFocus() {
    this.refs.filter.getInputDOMNode().select();
  }

  render() {
    const { channels, unfilteredChannels, isLoading } = this.props;

    if (isLoading) {
      return <Loading />;
    }

    if (_.isEmpty(unfilteredChannels) && !isLoading) {
      return (
        <span>You haven't subscribed to any channels yet.
          Discover new channels and podcasts <Link to="/search/">here</Link>.</span>);
    }

    return (
      <div>
        <Input className="form-control"
               type="search"
               ref="filter"
               onClick={this.handleFocus.bind(this)}
               onKeyUp={this.handleFilterChannels.bind(this)}
               placeholder="Find a channel" />
        <Input>
          <a className="btn btn-default form-control"
            href={`/podbaby-${moment().format('YYYY-MM-DD')}.opml`} download><Icon icon="download" /> Download OPML</a>
        </Input>
      {this.props.channels.map(channel => {
        const toggleSubscribe = () => {
            this.props.dispatch(actions.subscribe.toggleSubscribe(channel));
        };
        return <ChannelItem key={channel.id}
                            channel={channel}
                            isLoggedIn={true}
                            subscribe={toggleSubscribe} />;
      })}
      </div>
    );
  }
}

Subscriptions.propTypes = {
    channels: PropTypes.array.isRequired
};

const mapStateToProps = state => {
  return {
    isLoading: state.channels.isLoading,
    ...channelsSelector(state)
  };
};

export default connect(mapStateToProps)(Subscriptions);
