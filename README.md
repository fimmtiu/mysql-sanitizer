# mysql-sanitizer

This daemon sits between a MySQL client and server, transparently substituting values in any non-whitelisted string column with garbage. It's suitable for giving access to a production server while preventing users from seeing the entire contents of the database.

**In theory.**

In practice, this program was hacked together in about a day and a half by multiple people working as fast as they could with multiple false starts. The code in here is not production-ready and should not be taken as an example of how to do anything. Still, it seems to work.

Note that, since it's more of a proof-of-concept than a finished program, it presently only sanitizes responses from regular MySQL queries. Here are some things we haven't implemented or tested yet:

* Results from [executing a prepared statement](https://dev.mysql.com/doc/internals/en/com-stmt-execute.html)
* Results from [multi-statement queries](https://dev.mysql.com/doc/internals/en/multi-statement.html)
* Results from [executing stored procedures](https://dev.mysql.com/doc/internals/en/stored-procedures.html), including [multiple result sets](https://dev.mysql.com/doc/internals/en/multi-resultset.html)

## TODO

* Time out if a query has been in progress for >20 seconds. (Just kill the query, not the entire connection.)

* Consider removing mysqlproto entirely and rolling our own packet stuff. It's not great, and didn't buy us nearly as much as we'd hoped.

* Will our rewriting strings cause us to exceed `max_allowed_packet` for very large requests?
