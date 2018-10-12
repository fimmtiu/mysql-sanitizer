# mysql-sanitizer

This daemon sits between a MySQL client and server, transparently substituting values in any non-whitelisted string column with garbage. It's suitable for giving access to a production server while preventing users from seeing the entire contents of the database.

**In theory.**

In practice, this program was hacked together in about a day and a half by multiple people working as fast as they could with multiple false starts. The code in here is not production-ready and should not be taken as an example of how to do anything. Still, it seems to work.

Note that, since it's more of a proof-of-concept than a finished program, it presently only sanitizes responses from regular MySQL queries. Attempting to use any features that we don't currently handle ([prepared statements](https://dev.mysql.com/doc/internals/en/com-stmt-execute.html), [stored procedures](https://dev.mysql.com/doc/internals/en/stored-procedures.html), [multi-statement queries](https://dev.mysql.com/doc/internals/en/multi-statement.html), etc.) will signal an error.

We also currently don't allow returning the results of any MySQL function call; any string returned from a function will always be sanitized.

## Future work

* We need authentication, perhaps via Infrastructure credentials or Google SSO.

* We need to log all of the following: connection, disconnection, command executed, number of rows returned. This must include the username and IP of the user.

* We need to prevent people testing for the existence of records by doing something like `SELECT * FROM users WHERE first_name = "Bob" AND last_name = "Smith"`. Perhaps we can parse the SQL with something like [https://github.com/xwb1989/sqlparser](https://github.com/xwb1989/sqlparser), then barf if the `WHERE` clause contains anything we would sanitize? More investigation needed.

## TODO

* Consider removing mysqlproto entirely and rolling our own packet stuff. It's not great, and didn't buy us nearly as much as we'd hoped.

* ClientChannel and ServerChannel have terrible names. Fix that.

* ProxyConnection is kind of an orphan, and ServerConnection is much too large.

* No tests for the fragile state machine code. Need to get on this.
