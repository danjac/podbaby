import _ from 'lodash';
import React from 'react';
import TestUtils from 'react-addons-test-utils';
import jsdom from 'mocha-jsdom';
import { assert } from 'chai';

import { Podcast } from '../components/podcasts';

const makePodcast = attrs => {
  return {
    id: 1000,
    title: "test",
    channelId: 1000,
    name: "My Channel",
    ...attrs || {}
  };
};

const makePodcastProps = (podcast, props={}) => {
  return {
    podcast,
    togglePlayer: _.noop,
    toggleSubscribe: _.noop,
    toggleDetail: _.noop,
    toggleBookmark: _.noop,
    showChannel: true,
    isPlaying: false,
    channelUrl: "/channel/11/",
    ...props
  }

};

class Wrapper extends React.Component {
  render() {
    return (
      <div>{this.props.children}</div>
    );
  }
}

describe('Podcast component', function() {

  jsdom({ skipWindowCheck: true });

  it('should show channel if showChannel is true', function() {
    const podcast = makePodcast();
    const props = makePodcastProps(podcast, { showChannel: true });
    const component = <Wrapper><Podcast {...props} /></Wrapper>;
    const rendered = TestUtils.renderIntoDocument(component, 'div');
    assert(rendered, 'rendered is null');
const tags = TestUtils.scryRenderedDOMComponentsWithClass(rendered, "media-object")
    assert.equal(tags.length, 1)

//const shallowRenderer = TestUtils.createRenderer();
//    shallowRenderer.render(<Podcast {...props} />);
//    const result = shallowRenderer.getRenderOutput();

  });

  it('should not show channel if showChannel is false', function() {
    const podcast = makePodcast();
    const props = makePodcastProps(podcast, { showChannel: false });
    const component = <Wrapper><Podcast {...props} /></Wrapper>;
    const rendered = TestUtils.renderIntoDocument(component, 'div');
    assert(rendered, 'rendered is null');
const tags = TestUtils.scryRenderedDOMComponentsWithClass(rendered, "media-object")
    assert.equal(tags.length, 0)

//const shallowRenderer = TestUtils.createRenderer();
//    shallowRenderer.render(<Podcast {...props} />);
//    const result = shallowRenderer.getRenderOutput();

  });

});
