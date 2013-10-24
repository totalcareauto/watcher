# The Watcher

### Requirements

- Go (`brew install` go for macos with brew)
- Go cross compiliation support (see https://github.com/davecheney/golang-crosscompile)

### To Build

- cd to project root
- ``export GOPATH=`pwd` ``
- `GOOS=windows GOARCH=amd64 go install watcher`
- `GOOS=windows GOARCH=386 go install watcher`

### Example Config File

```json
{
  "watch": "/path/to/file/to/be/watched",
  "files": {
    "/path/to/file/to/be/uploaded":
    [
      "http://url.to/upload/to",
      "http://url.to/upload/to-stage"
    ],
    "/another/path/to/file/to/be/uploaded":
    [
      "http://another.url.to/upload/to",
      "http://another.url.to/upload/to-stage"
    ]
  }
}
```
