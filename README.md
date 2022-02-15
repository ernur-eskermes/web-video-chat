# Web Video Chat [Backend Application] ![GO][go-badge]

[go-badge]: https://img.shields.io/github/go-mod/go-version/p12s/furniture-store?style=plastic
[go-url]: https://github.com/p12s/furniture-store/blob/master/go.mod

Learn More about project [here](https://google.com)

## Build & Run (Locally)
### Prerequisites
- go 1.17
- docker & docker-compose
- [golangci-lint](https://github.com/golangci/golangci-lint) (<i>optional</i>, used to run code checks)
- [swag](https://github.com/swaggo/swag) (<i>optional</i>, used to re-generate swagger documentation)

Create .env file in root directory and add following values:
```dotenv
APP_ENV=local

MONGO_URI=mongodb://mongodb:27017
MONGO_USER=admin
MONGO_PASS=qwerty

PASSWORD_SALT=<random string>
JWT_SIGNING_KEY=<random string>

HTTP_HOST=localhost
```

Use `make run` to build&run project, `make lint` to check code with linter.