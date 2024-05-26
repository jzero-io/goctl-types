# goctl-types

goctl types group plugin

## install

```shell
go install github.com/jzero-io/goctl-types@latest
```

## Usage

```shell
goctl api plugin -plugin goctl-types="gen" -api main.api
```

## More Usage

```shell
goctl api plugin -plugin goctl-types="gen --filename-template={{.group}}.go" -api main.api --style go_zero
```