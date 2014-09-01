# Pavo

Server-side upload service for [jQuery-File-Upload](https://github.com/blueimp/jQuery-File-Upload).

## Install

[Install](https://golang.org/doc/install) Golang. Set the [GOPATH](http://golang.org/doc/code.html#GOPATH) environment variable. For example for MacOS:
```sh
brew install go
mkdir $HOME/go
# Add this line in your .zshrc or .bash_profile
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

Install application:
```sh
go get github.com/kavkaz/pavo
```

### Production: nginx setup

When used in a production environment it is recommended to use a web server nginx.

## License

[MIT license](http://www.opensource.org/licenses/MIT). Copyright (c) 2014 Zaur Abasmirzoev.

