import { createSelector } from 'reselect';

const subscriptionsSelector = state => state.subscriptions || [];
const relatedChannelsPreSelector = state => state.channel.relatedChannels;
const channelsPreSelector = state => state.channels.channels;
const channelPreSelector = state => state.channel.channel;
const filterSelector = state => state.channels.filter;
const currentPageSelector = state => state.channels.page;

const isSubscribed = (channel, subscriptions) => {
  return subscriptions.includes(channel.id);
};

export const channelSelector = createSelector(
  [channelPreSelector, subscriptionsSelector],
  (channel, subscriptions) => {
    if (!channel) {
      return null;
    }
    return Object.assign({}, channel, {
      isSubscribed: isSubscribed(channel, subscriptions),
    });
  }
);

export const relatedChannelsSelector = createSelector(
  [relatedChannelsPreSelector,
   subscriptionsSelector],
  (channels, subscriptions) => {
    return channels.map(channel => {
      return Object.assign({}, channel, { isSubscribed: isSubscribed(channel, subscriptions) });
    });
  });

export const channelsSelector = createSelector(
  [channelsPreSelector,
   filterSelector,
   subscriptionsSelector,
   currentPageSelector],
  (channels, filter, subscriptions, currentPage) => {
    const unfilteredChannels = channels.map(channel => {
      return Object.assign({}, channel, { isSubscribed: isSubscribed(channel, subscriptions) });
    });
    const filterToLower = filter.toLowerCase();

    const filteredChannels = unfilteredChannels.filter(channel => {
      return !filter || channel.title.toLowerCase().indexOf(filterToLower) > -1;
    });

    const pageSize = 10;

    const numPages = Math.ceil(filteredChannels.length / (pageSize * 1.0));
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
