import { createSelector } from 'reselect';

const subscriptionsSelector = state => state.subscriptions || [];
const channelsPreSelector = state => state.channels.channels;
const channelPreSelector = state => state.channel.channel;
const filterSelector = state => state.channels.filter;

const isSubscribed = (channel, subscriptions) => {
  return subscriptions.includes(channel.id);
};

export const channelSelector = createSelector(
  [ channelPreSelector, subscriptionsSelector ],
  (channel, subscriptions) => {
    if (!channel) {
      return null;
    }
    return Object.assign({}, channel, {
      isSubscribed: isSubscribed(channel, subscriptions)
    });
  }
);

export const channelsSelector = createSelector(
  [ channelsPreSelector,
    filterSelector,
    subscriptionsSelector ],
  (channels, filter, subscriptions) => {

    const unfilteredChannels = channels.map(channel => {
      channel.isSubscribed = isSubscribed(channel, subscriptions);
      return channel;
    });

    const filteredChannels = unfilteredChannels.filter(channel => {
      return !filter || channel.title.toLowerCase().indexOf(filter) > -1;
    });

    return {
      channels: filteredChannels,
      unfilteredChannels,
      filter
    };
  }

);
