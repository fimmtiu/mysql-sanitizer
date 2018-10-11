# mysql-sanitizer

Brightens, whitens, and shines dingy, dusty, and dirty data!

## TODO

* Get a list of whitelisted columns (JSON file?) and read it on startup.

* Write the sanitizer to remove non-whitelisted columns from responses.

* Store credentials in the proxy, then replace the credentials in the handshake packet with the stored credentials.
