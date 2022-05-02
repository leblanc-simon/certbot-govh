# certbot-govh

certbot-govh is a simple plugin for certbot to create and delete DNS entry for DNS Validation on OVH.

## Installation

```
cd /opt
git clone https://github.com/leblanc-simon/certbot-govh.git
cd certbot-govh
go build
```

Or download release : https://github.com/leblanc-simon/certbot-govh/releases/download/v1.0.0/certbot-govh-1.0.0.tar.bz2

## Usage

First, you need create a token for your domain. You may to directly go to : [https://eu.api.ovh.com/createToken/](https://eu.api.ovh.com/createToken/).

Set a reasonable validity, then set the right for your domain name :
* POST : /domain/zone/{domain name}/record
* POST : /domain/zone/{domain name}/refresh
* DELETE : /domain/zone/{domain name}/record/*

Restrict by IP if possible.

Create, a directory in your home directory named ".ovh", and then, a file in this directory named "{domain name}.ini" : 

```ini
# /root/.ovh/example.com.ini
[default]
endpoint="ovh-eu"

[ovh-eu]
application_key=
application_secret=
consumer_key=
```

When all is OK, you can run : 

```bash
certbot certonly --manual --preferred-challenges=dns --manual-auth-hook /opt/certbot-govh/certbot-govh-auth.sh --manual-cleanup-hook /opt/certbot-govh/certbot-govh-cleanup.sh -d *.example.com
```

Or use a config file : 

```ini
# /etc/letsencrypt/cli-example.com.ini
rsa-key-size = 4096
email = postmaster@example.cpù
agree-tos = True

domains = *.example.com

authenticator = manual
manual-auth-hook = /opt/certbot-govh/certbot-govh-auth.sh
manual-cleanup-hook = /opt/certbot-govh/certbot-govh-cleanup.sh
preferred-challenges = dns-01
manual-public-ip-logging-ok = True
renew-by-default = True
```

```bash
certbot certonly -c /etc/letsencrypt/cli-example.com.ini
```

## License

* [WTFPL](http://www.wtfpl.net/)

## Author

Simon Leblanc <contact@leblanc-simon.eu>