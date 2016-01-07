import { createSelector } from 'reselect';

const podcastsPreSelector = state => state.podcasts.podcasts || [];
const podcastPreSelector = state => state.podcast.podcast || [];

const playerSelector = state => state.player;
const detailSelector = state => state.podcasts.showDetail || [];
const bookmarksSelector = state => state.bookmarks.bookmarks || [];

const isBookmarked = (bookmarks, podcast) => {
  return bookmarks.includes(podcast.id);
};

const isShowDetail = (showDetail, podcast) => {
  return showDetail.includes(podcast.id);
};

const isPlaying = (player, podcast) => {
  return player.podcast && player.podcast.id === podcast.id;
};

const assign = (podcast, bookmarks, showDetail, player) => {
  if (!podcast) return null;
  return Object.assign({}, podcast, {
    isBookmarked: isBookmarked(bookmarks, podcast),
    isShowDetail: isShowDetail(showDetail, podcast),
    isPlaying: isPlaying(player, podcast)
  });
};

export const podcastSelector = createSelector(
  [ podcastPreSelector,
    detailSelector,
    bookmarksSelector,
    playerSelector ],
  (podcast, showDetail, bookmarks, player) => {
    return assign(podcast, bookmarks, showDetail, player);
  }
);

export const podcastsSelector = createSelector(
  [ podcastsPreSelector,
    detailSelector,
    bookmarksSelector,
    playerSelector ],
  (podcasts, showDetail, bookmarks, player) => {
    return podcasts.map(podcast => {
      return assign(podcast, bookmarks, showDetail, player);
    });
  }
);
