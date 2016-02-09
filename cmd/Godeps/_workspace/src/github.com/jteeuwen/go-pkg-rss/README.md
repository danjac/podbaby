# RSS

This package allows us to fetch Rss and Atom feeds from the internet.
They are parsed into an object tree which is a hybrid of both the RSS and Atom
standards.

Supported feeds are:
- Rss v0.91, 0.92 and 2.0
- Atom 1.0

The package allows us to maintain cache timeout management. This prevents us
from querying the servers for feed updates too often and risk ip bans. Apart 
from setting a cache timeout manually, the package also optionally adheres to
the TTL, SkipDays and SkipHours values specified in the feeds themselves.

Note that the TTL, SkipDays and SkipHour fields are only part of the RSS spec.
For Atom feeds, we use the CacheTimeout in the Feed struct.

Because the object structure is a hybrid between both RSS and Atom specs, not
all fields will be filled when requesting either an RSS or Atom feed. I have
tried to create as many shared fields as possible but some of them simply do
not occur in either the RSS or Atom spec.

The Feed object supports notifications of new channels and items.
This is achieved by passing 2 function handlers to the feeder.New() function.
They will be called whenever a feed is updated from the remote source and
either a new channel or a new item is found that previously did not exist.
This allows you to easily monitor a feed for changes. See feed_test.go for
an example of how this works.

## DEPENDENCIES

[github.com/jteeuwen/go-pkg-xmlx](http://github.com/jteeuwen/go-pkg-xmlx)

## USAGE


An idiomatic example program can be found in [testdata/example.go](https://github.com/jteeuwen/go-pkg-rss/blob/master/testdata/example.go).

