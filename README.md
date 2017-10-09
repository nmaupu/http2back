# What is http2back

http2back provides a HTTP server to upload files to various backends :

* Filesystem
* FTP
* AWS S3
* Dropbox

# Building

Requirements :

- Go lang > 1.5
- glide

```
git clone https://github.com/nmaupu/http2back
cd http2back
glide install && make
```

# Usage

Once configured (see next section), just run the server :

```
./http2back
2017/10/09 23:51:47 Starting http server on localhost:4242 using provider Filesystem (/tmp)
```

Send something using the client of your choice :

```
curl -X PUT -F file=@file-42.lol http://localhost:4242/
```


# Configuration

As for now, configuration can only be done using a config file named http2back.{json, toml, yaml, hcl, properties}.
This config file has to be present in http2back classpath including : /etc/http2back, $HOME/.http2back or in the current directory.
Have a look at the sample file at the root of the repository : `http2back.yaml`

Configuration is easy (using yaml for the win):

```
bind_address: 127.0.0.1
port: 4242
provider:
  name: <provider name>
  [more provider opts]
```

# Providers

## Filesystem

```
provider:
  name: filesystem
  dest: /tmp
```

## FTP

```
provider:
  name: ftp
  dest: /
  host: host:port
  username: user
  password: pass
```

## AWS S3

```
provider:
  name: s3
  dest: /
  bucket: my-bucket
  region: eu-west-1
  aws-access-key-id: my-access-key-id
  aws-secret-access-key: my-secret-access-key
```

## Dropbox

First, you need to create an app and an associated token.

```
provider:
  name: dropbox
  token: <my api token>
  dest: /myfiles
```

# Dependencies

The following dependencies have been used :

* https://github.com/jawher/mow.cli
* https://github.com/jlaffaye/ftp
* https://github.com/aws/aws-sdk-go
* https://github.com/tj/go-dropbox

