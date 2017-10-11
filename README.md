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
./http2back filesystem --dest=/tmp
2017/10/11 11:30:44 Starting http server on 127.0.0.1:8080 using provider Filesystem (tmp)
```

Send something using the client of your choice :

```
curl -X PUT -F file=@file-42.lol http://127.0.0.1:8080/
```


# Configuration

Configuration is done using command arguments :

```
./http2back [--bind=<binding-address>] [--port=<server-port>] COMMAND [OPTS]
```

Each command corresponds to one provider, pop help using :

```
./http2back COMMAND help
```

# Dependencies

The following dependencies have been used :

* https://github.com/jawher/mow.cli
* https://github.com/jlaffaye/ftp
* https://github.com/aws/aws-sdk-go
* https://github.com/tj/go-dropbox

