import React from 'react';
import TestUtils from 'react-addons-test-utils';
import { assert } from 'chai';
import jsdom from 'jsdom-global';

import Player from '../../components/player';
import { Wrapper } from './utils';
import { makePodcast, makePlayerProps } from './fixtures';

describe('Player component', function () {
  before(function () {
    this.jsdom = jsdom();
  });

  after(function () {
    this.jsdom();
  });

  it('should render the truncated podcast title', function () {
    const podcast = makePodcast({ name: 'We do cool podcasts', title: 'Some title' });
    const props = makePlayerProps(podcast);
    const component = <Wrapper><Player {...props} /></Wrapper>;
    const rendered = TestUtils.renderIntoDocument(component, 'div');

    const $title = TestUtils.findRenderedDOMComponentWithTag(rendered, 'b');
    const $link = $title.children[0];

    const title = $link.getAttribute('title');
    assert.equal(title, 'We do cool podcasts : Some title');
    assert.equal($link.textContent, 'We do cool podcasts : Some ...');
  });
});
