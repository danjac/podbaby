import _ from 'lodash';
import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import DocumentTitle from 'react-document-title';

import * as actions from '../actions';
import { channelsSelector } from '../selectors';
import Loading from '../components/loading';
import ChannelItem from '../components/channel_item';
import { getTitle } from './utils';

export class Recommendations extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindActionCreators(actions.channels, dispatch);
  }

  render() {
    const { channels, isLoading, isLoggedIn } = this.props;

    if (isLoading) {
      return <Loading />;
    }

    if (_.isEmpty(channels) && !isLoading) {
      return (
        <span>We can't find any recommendations for you at the moment.
          Discover other feeds and podcast episodes <Link to="/search/">here</Link>.</span>);
    }

    return (
      <DocumentTitle title={getTitle('Recommendations')}>
      <div>
      {this.props.channels.map(channel => {
        const toggleSubscribe = () => {
          this.props.dispatch(actions.subscribe.toggleSubscribe(channel));
        };
        return (
          <ChannelItem
            key={channel.id}
            channel={channel}
            isLoggedIn={isLoggedIn}
            subscribe={toggleSubscribe}
          />
        );
      })}
      </div>
    </DocumentTitle>
    );
  }
}

Recommendations.propTypes = {
  dispatch: PropTypes.func.isRequired,
  isLoggedIn: PropTypes.bool.isRequired,
  isLoading: PropTypes.bool.isRequired,
  channels: PropTypes.array.isRequired,
};

const mapStateToProps = state => {
  const { isLoading } = state.channels;
  const { isLoggedIn } = state.auth;
  const { channels } = channelsSelector(state);

  return {
    isLoggedIn,
    isLoading,
    channels,
  };
};

export default connect(mapStateToProps)(Recommendations);
