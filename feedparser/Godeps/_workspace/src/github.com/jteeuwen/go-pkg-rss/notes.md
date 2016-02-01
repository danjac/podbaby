## The problem with the current database solution

The purpose of the database is to ensure that the channel and item
handlers are only called once for each new channel and each new
item. It is clear that many users of go-pkg-rss are having problems
with their channel and item handlers being called multiple times for
the same items.

The current solution makes writing handlers very clean and safe as the
user can be sure that their handlers are only executed once for each
new item/channel.

### Use cases ###

The use cases where database shines is with regards to long running go
routines what watch a specific feed url and over the lifetime of the
program periodically check that feed url for updates. In that
situation, having a database prevents the duplication of items as
there is a high likelyhood that the refetch of the feed url will
contain items already processed by the item handler.

The benefits of this include:

1. Batteries included: If the user is creating a program that
   processes a set of feed urls that it repeatedly polls on an
   existing feed then the built in database provides a "batteries
   included" solution to prevent calling the users channel or item
   handlers unnecessarily, greatly simplifying development of those
   handlers.

2. Zero maintenance: The database is maintence free. There is nothing
   the developer needs to do to take advantage of it. At the same time
   the developer doesn't need to query it at any time nor do they need
   to purge items from it. They need not even know it exists.

### Problems

The problem with the current solution is that it doesn't scale. Time
for an example. I'm writing a program that will pull feed urls from a
queue system to be processed. In order to execute several feed fetches
at once it may run several hundred go routines, each fetching a feed
url from the queue, processing it via go-pkg-rss and the putting the
item urls into another queue for processing by another program. The
feed url job is then released to the queue which will then delay the
feed url from being processed for a set amount of time (usually the
lastupdate/cacheTimeout). As there are several thousand feed urls to
get through, I will be running my program on several servers, each
fetching feed urls from the queue via its several hundred go routines.

In order to prevent duplication of effort, as across several thousand
feed urls there is a very high likelyhood that items will be
duplicated across feeds, I record a hash for each item in
memcached. This provides a very quick and lightweight way of
determining if I have already fetched that article before, and
therefore do not need to fetch it again. This has the added benefit
that the cache can be shared across several servers all collectiong
feed urls and article urls as a centralised "database" (although I am
also leaning on the caching features of memcache to store the entry
for a limited time, allowing me to fetch an article again in the
future, in case of updates to the article).

In addition to this, I also check and catch errors raised by network
issues such as timeouts, unparsable urls, http error codes and
unparsable documents. For these I also store a hash in memcache for
each feed url, however this is an incrementing value, allowing me to
keep track or the number of retry attempts made to fetch that feed
url. After a certain threshold is met, and the feed url is still
failing, I mark the job as bad in the queue system, which prevents me
from constantly refetching a bad feed url.

The current database solution contributes the following issues:

1. The database is too simple: In the above example I need to track
   (and prevent) the number of fetches I make for article items. I
   also need to allow a number of retry attempts before preventing
   refetches.The current database is not sophisticated enough to
   handle these two different cases at once.

2. The database does expire or clean up entries: I expect my program
   to run for a very long time, processing many thousand feeds and
   even more article urls. The current implementation of the database
   is simple in that it continues to grow indefinately, consuming lots
   of memory.

3. The database replaces a job that is trivial to implement in a
   handler: The current database doesn't provide anything that
   couldn't be developed by a user of the package with ease.

4. The database doesn't provide a fine-grained enough key for
   duplicates: The current version uses a multitude of options for
   item keys and channels, all of which could very easily be falsly
   marked as duplicates. For example, the item title is not a very
   unique key, expecially when each item and channel each have unique
   key in the form of the url.

### Proposed solution

Looking across the stdlib provided with Go we can see and example
where similar concerns have been met with elegance. The net/http
package uses handlers, much like go-pkg-rss does, to off load
implementation complexity to outside of the core http package. It also
provides batteris included solutions to common problems that
developers may have with built in handlers such as a FileServer, a
NotFoundHandler, RedirectHandler, StripPrefix and TimeoutHandler. I
propose that the current database implementation be stripped from the
package and moved to a set of built in handlers.

The developer will then be provided with two powerful options:

1. Use the built in database handlers as part of a chain along with
   their own existing handlers to get the same functionality as
   currently provided by the existing database implementation.

2. Roll their own custom handlers which may be inserted into the
   handler chain. They even have the option of using the provided
   database handlers if they want.

This opens up exciting possibilities for developers and the future of
the go-pkg-rss package. We could add handlers for

- caching via memcache,
- a retry count handler for channels,
- a whitelist/blacklist handler based on item title,
- a filter handler that strips out items that have a date older than 3
  hours,
- etc.

What's nice about this approach is that we can recreate a functionally
identical approach that means this is backwards compatible with older
code that implements the old way of dealing with handlers and
duplicate channels and items.
