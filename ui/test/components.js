import _ from 'lodash';
import React from 'react';
import TestUtils from 'react-addons-test-utils';
import { assert } from 'chai';
import jsdom from 'jsdom-global';

import Player from '../components/player';
import Podcast from '../components/podcast_item';

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
    showExpanded: false,
    isLoggedIn: true,
    isPlaying: false,
    channelUrl: "/channel/11/",
    ...props
  }

};

const makePlayerProps = (podcast, props={}) => {
  return {
    onClose: _.noop,
    onTimeUpdate: _.noop,
    onToggleBookmark: _.noop,
    isLoggedIn: true,
    player: {
      podcast,
      isPlaying: true,
    },
    ...props
  };
};

class Wrapper extends React.Component {
  render() {
    return (
      <div>{this.props.children}</div>
    );
  }
}

describe('Player component', function() {

  before(function() {
    this.jsdom = jsdom();
  });

  after(function() {
    this.jsdom();
  });

  it('should render the truncated podcast title', function() {

    const podcast = makePodcast({ name: 'We do cool podcasts', title: 'Some title' });
    const props = makePlayerProps(podcast);
    const component = <Wrapper><Player {...props} /></Wrapper>;
    const rendered = TestUtils.renderIntoDocument(component, 'div');

    const $title = TestUtils.findRenderedDOMComponentWithTag(rendered, "b");
    const $link = $title.children[0];

    const title = $link.getAttribute("title");
    assert.equal(title, "We do cool podcasts : Some title");
    assert.equal($link.textContent, "We do cool podcasts : Some ...");

  });

});

describe('Podcast component', function() {

  before(function() {
    this.jsdom = jsdom();
  });

  after(function() {
    this.jsdom();
  });

  it('should show remove bookmark button if is bookmarked', function() {

    const podcast = makePodcast({ isBookmarked: true });
    const props = makePodcastProps(podcast);
    const component = <Wrapper><Podcast {...props} /></Wrapper>;
    const rendered = TestUtils.renderIntoDocument(component, 'div');
    const $buttons = TestUtils.scryRenderedDOMComponentsWithTag(rendered, 'button');

    const titles = $buttons.map(node => node.getAttribute("title"));
    assert.include(titles, 'Remove bookmark');
  });

  it('should show bookmark button if is bookmarked', function() {

    const podcast = makePodcast({ isBookmarked: false });
    const props = makePodcastProps(podcast);
    const component = <Wrapper><Podcast {...props} /></Wrapper>;
    const rendered = TestUtils.renderIntoDocument(component, 'div');
    const $buttons = TestUtils.scryRenderedDOMComponentsWithTag(rendered, 'button');

    const titles = $buttons.map(node => node.getAttribute("title"));
    assert.include(titles, 'Add to bookmarks');
  });

  it('should show channel if showChannel is true', function() {
    const podcast = makePodcast();
    const props = makePodcastProps(podcast, { showChannel: true });
    const component = <Wrapper><Podcast {...props} /></Wrapper>;
    const rendered = TestUtils.renderIntoDocument(component, 'div');
    const tags = TestUtils.scryRenderedDOMComponentsWithClass(rendered, "media-body")
    assert.equal(tags.length, 1)
    const $header = TestUtils.findRenderedDOMComponentWithTag(rendered, 'h5');
    assert.equal($header.textContent, podcast.name);

  });

  it('should not show channel if showChannel is false', function() {
    const podcast = makePodcast();
    const props = makePodcastProps(podcast, { showChannel: false });
    const component = <Wrapper><Podcast {...props} /></Wrapper>;
    const rendered = TestUtils.renderIntoDocument(component, 'div');
    const tags = TestUtils.scryRenderedDOMComponentsWithClass(rendered, "media-object")
    assert.equal(tags.length, 0)

//const shallowRenderer = TestUtils.createRenderer();
//    shallowRenderer.render(<Podcast {...props} />);
//    const result = shallowRenderer.getRenderOutput();

  });

});
