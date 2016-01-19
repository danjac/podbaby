import immutable from 'immutable';

export default immutable.Record({
  id: null,
  image: '',
  title: '',
  isSubscribed: false,
  description: '',
  url: '',
  website: immutable.Map({
    Valid: false,
    String: '',
  }),
});
