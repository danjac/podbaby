import { createSelector } from 'reselect';

const subscriptionsSelector = state => state.subscriptions || [];
const channelsSelector = state => state.channels.channels;
const filterSelector = state => state.channels.filter;

export default createSelector(
  [ channelsSelector,
    filterSelector,
    subscriptionsSelector ],
  (channels, filter, subscriptions) => {

    const unfilteredChannels = channels.map(channel => {
      channel.isSubscribed = subscriptions.includes(channel.id);
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
