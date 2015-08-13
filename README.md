# tropo-http-auth
Simple mod_authnz_external app to authenticate users against Tropo Provisioning


### Usage

Usage is simple

* Install package ( `yum install mod_authnz_external` )
* Load module in httpd.conf
```
LoadModule authnz_external_module modules/mod_authnz_external.so

<VirtualHost *:80>
  ServerName        dashing.tropo.com

  <Location />
    AddExternalAuth tropo /usr/sbin/tropo-http-auth
    SetExternalAuthMethod tropo pipe
    RequestHeader    unset Accept-Encoding
    Order allow,deny
    Allow from all
  </Location>
</VirtualHost>
```
