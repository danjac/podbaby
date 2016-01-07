import { createSelector } from 'reselect';

const playerSelector = state => state.player;
const podcastsSelector = state => state.podcasts.podcasts || [];
const detailSelector = state => state.podcasts.showDetail || [];
const bookmarksSelector = state => state.bookmarks.bookmarks || [];

export default createSelector(
  [ podcastsSelector,
    detailSelector,
    bookmarksSelector,
    playerSelector ],
  (podcasts, showDetail, bookmarks, player) => {
    return podcasts.map(podcast => {
      podcast.isBookmarked = bookmarks.includes(podcast.id);
      podcast.isShowDetail = showDetail.includes(podcast.id);
      podcast.isPlaying = player.podcast && player.podcast.id === podcast.id;
      return podcast;
    });
  }
);
