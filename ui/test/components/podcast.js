import React from 'react';
import TestUtils from 'react-addons-test-utils';
import { assert } from 'chai';
import jsdom from 'jsdom-global';

import Podcast from '../../components/podcast_item';
import { Wrapper } from './utils';
import { makePodcast, makePodcastProps } from './fixtures';


describe('Podcast component', function () {
  before(function () {
    this.jsdom = jsdom();
  });

  after(function () {
    this.jsdom();
  });

  it('should show remove bookmark button if is bookmarked', function () {
    const podcast = makePodcast({ isBookmarked: true });
    const props = makePodcastProps(podcast);
    const component = <Wrapper><Podcast {...props} /></Wrapper>;
    const rendered = TestUtils.renderIntoDocument(component, 'div');
    const $buttons = TestUtils.scryRenderedDOMComponentsWithTag(rendered, 'button');

    const titles = $buttons.map(node => node.getAttribute('title'));
    assert.include(titles, 'Remove bookmark');
  });

  it('should show bookmark button if is bookmarked', function () {
    const podcast = makePodcast({ isBookmarked: false });
    const props = makePodcastProps(podcast);
    const component = <Wrapper><Podcast {...props} /></Wrapper>;
    const rendered = TestUtils.renderIntoDocument(component, 'div');
    const $buttons = TestUtils.scryRenderedDOMComponentsWithTag(rendered, 'button');

    const titles = $buttons.map(node => node.getAttribute('title'));
    assert.include(titles, 'Add to bookmarks');
  });

  it('should show channel if showChannel is true', function () {
    const podcast = makePodcast();
    const props = makePodcastProps(podcast, { showChannel: true });
    const component = <Wrapper><Podcast {...props} /></Wrapper>;
    const rendered = TestUtils.renderIntoDocument(component, 'div');
    const tags = TestUtils.scryRenderedDOMComponentsWithClass(rendered, 'media-body');
    assert.equal(tags.length, 1);
    const $header = TestUtils.findRenderedDOMComponentWithTag(rendered, 'h5');
    assert.equal($header.textContent, podcast.name);
  });

  it('should not show channel if showChannel is false', function () {
    const podcast = makePodcast();
    const props = makePodcastProps(podcast, { showChannel: false });
    const component = <Wrapper><Podcast {...props} /></Wrapper>;
    const rendered = TestUtils.renderIntoDocument(component, 'div');
    const tags = TestUtils.scryRenderedDOMComponentsWithClass(rendered, 'media-object');
    assert.equal(tags.length, 0);
  });
});
