[![CircleCI](https://circleci.com/gh/thoeni/go-tfl.svg?style=svg)](https://circleci.com/gh/thoeni/go-tfl)[![Coverage Status](https://coveralls.io/repos/github/thoeni/go-tfl/badge.svg?branch=master)](https://coveralls.io/github/thoeni/go-tfl?branch=master)[![Go Report Card](https://goreportcard.com/badge/github.com/thoeni/go-tfl)](https://goreportcard.com/report/github.com/thoeni/go-tfl)

##What
go-tfl is a go library that provides a client to be used to query TFL APIs.

The current status of the project is in its early stages, therefore many things will change and for now it's used by the [slack-tube-service](https://github.com/thoeni/slack-tube-service/tree/refactor/tflclient-package)

##How do I use it?

Create a tfl Client, call the `GetTubeStatus` method and you'll get back a `[]Report` array

```golang
import "github.com/thoeni/go-tfl"

func main() {
  client := tfl.NewClient()
  reports, err := client.GetTubeStatus()
  if err != nil {
    log.Print("Error while retrieving Tube statuses")
    return err
  }
}
```
