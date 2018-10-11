# mysql-sanitizer

Brightens, whitens, and shines dingy, dusty, and dirty data!

## TODO

* Make the JSON whitelist file a config variable which the user can specify.

* Write the sanitizer to remove non-whitelisted columns from responses.

* Store credentials in the proxy, then replace the credentials in the handshake packet with the stored credentials.

* When a connection exits, mysql-sanitizer starts consuming 200% CPU like crazy? WTF?
