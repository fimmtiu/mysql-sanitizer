# mysql-sanitizer

Brightens, whitens, and shines dingy, dusty, and dirty data!

## TODO

* Write the sanitizer to remove non-whitelisted columns from responses.

* Time out if a query has been in progress for >20 seconds. (Just kill the query, not the entire connection.)

* Consider removing mysqlproto entirely and rolling our own packet stuff. It's not great.

* Will our rewriting strings cause us to exceed `max_allowed_packet` for very large requests?

* We should add support for prepared statement responses, too.
