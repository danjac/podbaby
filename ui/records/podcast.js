import immutable from 'immutable';

export default immutable.Record({
  id: null,
  image: '',
  title: '',
  name: '',
  channelId: null,
  enclosureUrl: '',
  isPlaying: false,
  isBookmarked: false,
  isShowDetail: false,
  pubDate: new Date(),
  source: '',
  description: '',
});
