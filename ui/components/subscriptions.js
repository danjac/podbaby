import _ from 'lodash';
import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import moment from 'moment';

import {
  Grid,
  Row,
  Col,
  Button,
  ButtonGroup,
  Well,
  Panel,
  Input
} from 'react-bootstrap';

import * as  actions from '../actions';
import { channelsSelector } from '../selectors';
import Image from './image';
import Icon from './icon';
import Loading from './loading';

const ListItem = props => {
  const { channel, toggleSubscribe } = props;
  return (
    <Panel>
    <div className="media">
      <div className="media-left">
        <a href="#">

        <Image className="media-object"
               src={channel.image}
               errSrc='/static/podcast.png'
               imgProps={{
               height:60,
               width:60,
               alt:channel.title }} />
        </a>
      </div>
      <div className="media-body">
        <Grid>
          <Row>
            <Col xs={6} md={9}>
              <h4 className="media-heading"><Link to={`/podcasts/channel/${channel.id}/`}>{channel.title}</Link></h4>
            </Col>
            <Col xs={6} md={3}>
              <ButtonGroup>
                <Button title={channel.isSubscribed ? 'Unsubscribe': 'Subscribe'}
                        onClick={toggleSubscribe}><Icon icon={channel.isSubscribed ? 'unlink' : 'link'} /> {channel.isSubscribed ? 'Unsubscribe': 'Subscribe'}</Button>
              </ButtonGroup>
            </Col>
          </Row>
        </Grid>
      </div>
    </div>
  </Panel>

  );
};


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
          Discover new channels and podcasts <Link to="/podcasts/search/">here</Link>.</span>);
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
        return <ListItem key={channel.id}
                         channel={channel}
                         toggleSubscribe={toggleSubscribe} />;
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
