import { createSelector } from 'reselect';

const podcastsPreSelector = state => state.podcasts.get('podcasts');
const podcastPreSelector = state => state.podcast.get('podcast');

const playerSelector = state => state.player;
const detailSelector = state => state.podcasts.get('showDetail');
const bookmarksSelector = state => state.bookmarks.get('bookmarks');

const isBookmarked = (bookmarks, podcast) => {
  return bookmarks.includes(podcast.get('id'));
};

const isShowDetail = (showDetail, podcast) => {
  return showDetail.includes(podcast.get('id'));
};

const isPlaying = (player, podcast) => {
  return player.get('podcast') === podcast;
};

const assign = (podcast, bookmarks, showDetail, player) => {
  if (!podcast) return null;
  return podcast
    .set('isBookmarked', isBookmarked(bookmarks, podcast))
    .set('isShowDetail', isShowDetail(showDetail, podcast))
    .set('isPlaying', isPlaying(player, podcast));
};

export const podcastSelector = createSelector(
  [podcastPreSelector,
   detailSelector,
   bookmarksSelector,
   playerSelector],
  (podcast, showDetail, bookmarks, player) => {
    return assign(podcast, bookmarks, showDetail, player);
  }
);

export const podcastsSelector = createSelector(
  [podcastsPreSelector,
   detailSelector,
   bookmarksSelector,
   playerSelector],
  (podcasts, showDetail, bookmarks, player) => {
    return podcasts.map(podcast => {
      return assign(podcast, bookmarks, showDetail, player);
    });
  }
);
