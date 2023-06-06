# Changelog

## 1.0.0

### Features

* `admin` protocol for dedicated administrator connections

### Changed

* Added `Hidden()` method to `ProtocolParser` interface

## 0.21.0

### Features

* Updated azidentity to 1.2.1, which adds in memory cache for managed credentials ([#90](https://github.com/microsoft/go-mssqldb/pull/90))

### Bug fixes

* Fixed uninitialized server name in TLS config ([#93](https://github.com/microsoft/go-mssqldb/issues/93))([#94](https://github.com/microsoft/go-mssqldb/pull/94))
* Fixed several kerberos authentication usages on Linux with new krb5 authentication provider. ([#65](https://github.com/microsoft/go-mssqldb/pull/65))

### Changed

* New kerberos authenticator implementation uses more explicit connection string parameters.

| Old          | New                |
|--------------|--------------------|
| krb5conffile | krb5-configfile    |
| krbcache     | krb5-credcachefile |
| keytabfile   | krb5-keytabfile    |
| realm        | krb5-realm         |

## 0.20.0

### Features

* Add driver version and name to TDS login packets
* Add `pipe` connection string parameter for named pipe dialer
* Expose network errors that occur during connection establishment. Now they are
wrapped, and can be detected by using errors.As/Is practise. This connection
errors can, and could even before, happen anytime the sql.DB doesn't have free
connection for executed query.

### Bug fixes

* Added checks while reading prelogin for invalid data ([#64](https://github.com/microsoft/go-mssqldb/issues/64))([86ecefd8b](https://github.com/microsoft/go-mssqldb/commit/86ecefd8b57683aeb5ad9328066ee73fbccd62f5))

* Fixed multi-protocol dialer path to avoid unneeded SQL Browser queries

