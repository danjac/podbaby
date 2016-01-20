import _ from 'lodash';
import { createSelector } from 'reselect';

const podcastsPreSelector = state => state.podcasts.podcasts || [];
const podcastPreSelector = state => state.podcast.podcast || [];

const playerSelector = state => state.player;
const detailSelector = state => state.podcasts.showDetail || [];
const bookmarksSelector = state => state.bookmarks.bookmarks || [];
const playsSelector = state => state.plays || [];

const isBookmarked = (bookmarks, podcast) => {
  return bookmarks.includes(podcast.id);
};

const isShowDetail = (showDetail, podcast) => {
  return showDetail.includes(podcast.id);
};

const isPlaying = (player, podcast) => {
  return player.podcast && player.podcast.id === podcast.id;
};

const lastPlayedAt = (plays, podcast) => {
  const play = _.find(plays, value => value.podcastId === podcast.id);
  return play ? play.createdAt : null;
};

const assign = (podcast, bookmarks, showDetail, player, plays) => {
  if (!podcast) return null;
  return Object.assign({}, podcast, {
    isBookmarked: isBookmarked(bookmarks, podcast),
    isShowDetail: isShowDetail(showDetail, podcast),
    isPlaying: isPlaying(player, podcast),
    lastPlayedAt: lastPlayedAt(plays, podcast),
  });
};

export const podcastSelector = createSelector(
  [podcastPreSelector,
   detailSelector,
   bookmarksSelector,
   playerSelector,
   playsSelector,
  ],
  (podcast, showDetail, bookmarks, player, plays) => {
    return assign(podcast, bookmarks, showDetail, player, plays);
  }
);

export const podcastsSelector = createSelector(
  [podcastsPreSelector,
   detailSelector,
   bookmarksSelector,
   playerSelector,
   playsSelector,
  ],
  (podcasts, showDetail, bookmarks, player, plays) => {
    return podcasts.map(podcast => {
      return assign(podcast, bookmarks, showDetail, player, plays);
    });
  }
);
