import { Actions } from '../constants';

import { createAction } from './utils';

export function toggleDetail(podcast, isShowingDetail) {
  return isShowingDetail ? hidePodcastDetail(podcast.id) : showPodcastDetail(podcast.id);
}

export function hidePodcastDetail(podcastId) {
  return createAction(Actions.HIDE_PODCAST_DETAIL, podcastId);
}

export function showPodcastDetail(podcastId) {
  return createAction(Actions.SHOW_PODCAST_DETAIL, podcastId);
}
