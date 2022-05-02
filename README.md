# certbot-govh

certbot-govh is a simple plugin for certbot to create and delete DNS entry for DNS Validation on OVH.

## Installation

```
cd /opt
git clone https://github.com/leblanc-simon/certbot-govh.git
cd certbot-govh
go build
```

Or download release : 

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
certbot certonly --manual --preferred-challenges=dns --manual-auth-hook /opt/certbot-govh/certbot-govh-auth.sh --manual-cleanup-hook /opt/certbot-govh/certbot-govh-cleanup.sh -d secure.example.com
```

## License

* [WTFPL](http://www.wtfpl.net/)

## Author

Simon Leblanc <contact@leblanc-simon.eu>