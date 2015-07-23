## khabar-admin

Manage khabar notifications and settings

## Development

Clone the repo

```sh
$ go get github.com/codegangsta/gin
```

[gin](http://github.com/codegangsta/gin) is used to to automatically compile files while you are developing

Then run

```sh
$ go get && go install && PORT=7000 DEBUG=* gin -p 9000 -a 7000 -i run
```

Then visit `localhost:9000`

MongoDB config is stored in `middlewares/connect.go`.

[godep](https://github.com/tools/godep) is used for dependency management. So if you add or remove deps, make sure you run `godep save` before pushing code. Refer to its documentation for more info on how to use it.

## Usage

```sh
$ go get github.com/bulletind/khabar-admin
$ PORT=7000 khabar-admin # should start listening on port 7000
```

If you specify `MONGODB_URL` env variable then it will connect to that particular connection string.

#### Credits

Thanks to [go-martini](https://github.com/go-martini/martini), [go-mgo](https://github.com/go-mgo/mgo) and [godep](https://github.com/tools/godep)

## Todo

- Use [mux](www.gorillatoolkit.org/pkg/mux) instead of martini.
