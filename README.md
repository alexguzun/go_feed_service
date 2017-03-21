# go_feed_service
A service that reads, at a configured periodicity, RSS or Atom feeds and stores the items into a MongoDB database.
[![Build Status](https://travis-ci.org/alexguzun/go_feed_service.svg?branch=master)](https://travis-ci.org/alexguzun/go_feed_service)

##Database Model:

######Feed Source
```json
{
    "_id": "",
    "rss_urls": [],
    "description": "",
    "active": "",
    "use_feed_guid": "",
    "last_check": ""
}
```
######Feed entry
```json
{
    "_id": "",
    "title": "",
    "source_id": "",
    "content": "",
    "added_on": ""
}
```
##Configuration
The service should be configured with 2 environment variables:
- **FEED_READ_INTERVAL** defines the periodicity of feed read. Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h". If not set, has a default value of 5 Minute 
- **MONGO_URI** defines the MongoDB connection string. See [https://docs.mongodb.com/manual/reference/connection-string/](https://docs.mongodb.com/manual/reference/connection-string/) for more details.