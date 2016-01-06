import _ from 'lodash';
import React, { PropTypes } from 'react';
import { Link } from 'react-router';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

import {
  Grid,
  Row,
  Col,
  ButtonGroup,
  Button,
  Well,
  Input,
  Panel
} from 'react-bootstrap';

import  * as actions from '../actions';
import PodcastList from './podcasts';
import Image from './image';
import Icon from './icon';
import { sanitize, formatPubDate } from './utils';

const ChannelItem = props => {
  const { channel, subscribe } = props;
  const url = `/podcasts/channel/${channel.id}/`;

  return (
    <Panel>
    <div className="media">
      <div className="media-left">
        <Link to={url}>
        <Image className="media-object"
               src={channel.image}
               errSrc='/static/podcast.png'
               imgProps={{
               height:60,
               width:60,
               alt:channel.title }} />
        </Link>
      </div>
      <div className="media-body">
        <Grid>
          <Row>
            <Col xs={6} md={9}>
              <h4 className="media-heading"><Link to={url}>{channel.title}</Link></h4>
            </Col>
            <Col xs={6} md={3}>
              <ButtonGroup>
                <Button title={channel.isSubscribed ? "Unsubscribe" : "Subscribe"} onClick={subscribe}>
                  <Icon icon={channel.isSubscribed ? "unlink" : "link"} /> {channel.isSubscribed ? 'Unsubscribe' : 'Subscribe'}
                </Button>
              </ButtonGroup>
            </Col>
          </Row>
        </Grid>
      </div>
    </div>
  </Panel>
  );
};


export class Search extends React.Component {

  constructor(props) {
    super(props);
    const { search } = bindActionCreators(actions.search, this.props.dispatch);
    this.search = search;
  }

  componentDidMount() {
    const query = this.props.location.query.q || "";
    this.search(query);
    this.refs.query.getInputDOMNode().focus();
  }

  handleSearch(event) {
    event.preventDefault();
    const value = this.refs.query.getValue();
    this.search(_.trim(value));
  }

  handleFocus(event) {
    this.refs.query.getInputDOMNode().select();
  }

  render() {

    const { dispatch, channels, podcasts, searchQuery } = this.props;

    const help = (
      searchQuery ? '' :
        <span>
          <b>Hint:</b> Try a general category e.g. <em>history</em> or <em>movies</em>, the title of a podcast, or the name of a channel e. g. <em>RadioLab</em>.
        </span>
      );

    return (
    <div>
      <form className="form" onSubmit={this.handleSearch.bind(this)}>
        <Input type="search"
               ref="query"
               help={help}
               onClick={this.handleFocus.bind(this)}
               placeholder="Find a channel or podcast" />
        <Button type="submit" bsStyle="primary" className="form-control">
          <Icon icon="search" /> Search
        </Button>
      </form>
      {channels.map(channel => {
        const subscribe = (event) => {
          event.preventDefault();
          dispatch(actions.subscribe.toggleSubscribe(channel));
        };
        return (
          <ChannelItem key={channel.id}
                       channel={channel}
                       subscribe={subscribe} />
        );
      })}
      {podcasts.length > 0 ? <hr /> : ''}
      {searchQuery ?
        <PodcastList actions={actions}
                     showChannel={true}
                     ifEmpty=''
                      {...this.props} /> : '' }
    </div>
    );
  }
}

const mapStateToProps = state => {
  const { podcasts, showDetail, isLoading } = state.podcasts;
  const { query, channels } = state.search;
  return {
    searchQuery: query,
    podcasts: podcasts || [],
    channels: channels || [],
    showDetail,
    isLoading,
    player: state.player
  };
};

export default connect(mapStateToProps)(Search);
