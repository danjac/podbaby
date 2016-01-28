import _ from 'lodash';

export const makePodcast = (attrs = {}) => {
  return Object.assign({}, {
    id: 1000,
    title: 'test',
    channelId: 1000,
    name: 'My Channel',
    image: 'test.jpg',
  }, attrs);
};

export const makePodcastProps = (podcast, props = {}) => {
  return Object.assign({}, {
    podcast,
    togglePlayer: _.noop,
    toggleSubscribe: _.noop,
    toggleDetail: _.noop,
    toggleBookmark: _.noop,
    showChannel: true,
    showExpanded: false,
    isLoggedIn: true,
    isPlaying: false,
    channelUrl: '/channel/11/',
  }, props);
};

export const makePlayerProps = (podcast, props = {}) => {
  return Object.assign({}, {
    onClose: _.noop,
    onTimeUpdate: _.noop,
    onToggleBookmark: _.noop,
    onPlayNext: _.noop,
    onPlayLast: _.noop,
    isLoggedIn: true,
    player: {
      podcast,
      isPlaying: true,
    },
  }, props);
};
