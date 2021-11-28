## go-flickr

A minimalistic Flickr API client for Go

## Install

```
$ go get github.com/azer/go-flickr
```

## Manual

[API Reference](https://godoc.org/github.com/azer/go-flickr)

### Request

```go
import (
  "github.com/azer/go-flickr"
)

client := &flickr.Client{
  Key: "key",
  Token: "token", // optional
  Sig: "sig", // optional
}

response, err := client.Get("people.findByUsername", &flickr.Params{ "username": "azer" })
// => {"user":{"id":"98269877@N00", "nsid":"98269877@N00", "username":{"_content":"azerbike"}}, "stat":"ok"}
```

### FindUser

Find user by name.

```go
user, err := client.FindUser("azer")

user.Id
// => "123124324"

user.Name
// => azer
```

### Following

List the people given user follows on Flickr

```go
userId := "123123123"

following, err := client.Following(userId)
```

### Album

List photos in the album with given ID

```go
photos, err := client.Album("72157662053417706")
```

### Favorites

List photos that have been favorited by the given user ID. Note the Flickr API returns the list as paginated results.

```go
client, err := flickr.NewPhotosClient()
favs, err := client.Favs("98269877@N00")
if favs.Pages > 1 {
  favs, err = client.NextPage()
}
```
