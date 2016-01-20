import { createSelector } from 'reselect';

const subscriptionsSelector = state => state.subscriptions;
const channelsPreSelector = state => state.channels.get('channels');
const channelPreSelector = state => state.channel.get('channel');
const filterSelector = state => state.channels.get('filter');
const currentPageSelector = state => state.channels.get('page');

const isSubscribed = (channel, subscriptions) => {
  return subscriptions.includes(channel.get('id'));
};

export const channelSelector = createSelector(
  [channelPreSelector,
   subscriptionsSelector],
  (channel, subscriptions) => {
    if (!channel) {
      return null;
    }
    return channel.set('isSubscribed', isSubscribed(channel, subscriptions));
  }
);

export const channelsSelector = createSelector(
  [channelsPreSelector,
   filterSelector,
   subscriptionsSelector,
   currentPageSelector],
  (channels, filter, subscriptions, currentPage) => {
    const unfilteredChannels = channels.map(channel => {
      return channel.set('isSubscribed', isSubscribed(channel, subscriptions));
    });

    const filteredChannels = unfilteredChannels.filter(channel => {
      return !filter || channel.title.toLowerCase().indexOf(filter) > -1;
    });

    const pageSize = 10;

    const numPages = Math.ceil(filteredChannels.size / (pageSize * 1.0));
    const start = (currentPage - 1) * pageSize;
    const end = start + pageSize;

    const page = {
      page: currentPage,
      numPages,
    };

    const paginatedChannels = filteredChannels.slice(start, end);

    return {
      channels: paginatedChannels,
      unfilteredChannels,
      filter,
      page,
    };
  }

);
