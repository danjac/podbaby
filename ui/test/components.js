import _ from 'lodash';
import React from 'react';
import TestUtils from 'react-addons-test-utils';
import jsdom from 'mocha-jsdom';
import { assert } from 'chai';

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

class Wrapper extends React.Component {
  render() {
    return (
      <div>{this.props.children}</div>
    );
  }
}

describe('Podcast component', function() {

  jsdom({ skipWindowCheck: true });

  it('should show remove bookmark button if is bookmarked', function() {

    const podcast = makePodcast({ isBookmarked: true });
    const props = makePodcastProps(podcast);
    const component = <Wrapper><Podcast {...props} /></Wrapper>;
    const rendered = TestUtils.renderIntoDocument(component, 'div');
    const buttons = TestUtils.scryRenderedDOMComponentsWithTag(rendered, 'button');

    const titles = buttons.map(node => node.getAttribute("title"));
    assert.include(titles, 'Remove bookmark');
  });

  it('should show bookmark button if is bookmarked', function() {

    const podcast = makePodcast({ isBookmarked: false });
    const props = makePodcastProps(podcast);
    const component = <Wrapper><Podcast {...props} /></Wrapper>;
    const rendered = TestUtils.renderIntoDocument(component, 'div');
    const buttons = TestUtils.scryRenderedDOMComponentsWithTag(rendered, 'button');

    const titles = buttons.map(node => node.getAttribute("title"));
    assert.include(titles, 'Add to bookmarks');
  });

  it('should show channel if showChannel is true', function() {
    const podcast = makePodcast();
    const props = makePodcastProps(podcast, { showChannel: true });
    const component = <Wrapper><Podcast {...props} /></Wrapper>;
    const rendered = TestUtils.renderIntoDocument(component, 'div');
    const tags = TestUtils.scryRenderedDOMComponentsWithClass(rendered, "media-body")
    assert.equal(tags.length, 1)
    const h5 = TestUtils.findRenderedDOMComponentWithTag(rendered, 'h5');
    assert.equal(h5.textContent, podcast.name);

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
