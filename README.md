# Pavo

Server-side upload service for [jQuery-File-Upload](https://github.com/blueimp/jQuery-File-Upload) written in Golang.

## Usage

Running from the console by using the command:
```sh
$ pavo --host=localhost:9078 --storage=/path/to/root/storage
```

## Example

After install run application:
```sh
$ pavo --storage=$GOPATH/src/github.com/kavkaz/pavo/dummy/root_storage
```

Open example page in your browser:
```sh
open http://localhost:9073/example/jfu-basic.html
```

## Install

#### Install golang

[Install](https://golang.org/doc/install) Golang. Set the [GOPATH](http://golang.org/doc/code.html#GOPATH) environment variable. For example for MacOS:
```sh
brew install go
mkdir $HOME/go
# Add this line in your .zshrc or .bash_profile
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

#### Install application:

For first install run:
```sh
go get github.com/kavkaz/pavo
```

For update:
```sh
go get -u github.com/kavkaz/pavo
```

#### Setup nginx

When used in a production environment it is recommended to use a web server nginx. Configure the web server is reduced to specifying a directory for distribution static, location for the files, and optional authentication.

```
server {
    listen 80;
    server_name pavo.local;
    
    access_log /usr/local/var/log/nginx/pavo/access.log;
    error_log /usr/local/var/log/nginx/pavo/error.log notice;
    
    location /auth {
        internal;
        proxy_method GET;
        proxy_set_header Content-Length "";
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_pass_request_body off;
        proxy_pass http://localhost:3000/auth/url/in/your/app;
        proxy_set_header X-Original-URI $request_uri;
        client_max_body_size 0;
    }

    location /files {
        auth_request /auth;
    
        client_body_temp_path     /tmp;
        client_body_in_file_only  on;
        client_body_buffer_size   521K;
        client_max_body_size      10G;
    
        proxy_pass_request_headers on;
        proxy_set_header X-FILE $request_body_file;
        proxy_pass_request_body off;
        proxy_set_header Content-Length 0;
        proxy_pass http://127.0.0.1:9073;
    }
    
    location / {
        root /Path/To/Root/Of/Storage;
    }
}
```

These settings allow you to save the request body into a temporary file and pass on our application link to the file in the header `X-File`.

## License

[MIT license](http://www.opensource.org/licenses/MIT). Copyright (c) 2014 Zaur Abasmirzoev.

