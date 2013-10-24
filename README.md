# The Watcher

### Requirements

- Go (`brew install` go for macos with brew)
- Go cross compiliation support (see https://github.com/davecheney/golang-crosscompile)

### To Build

- cd to project root
- ``export GOPATH=`pwd` ``
- `GOOS=windows GOARCH=amd64 go install watcher`
- `GOOS=windows GOARCH=386 go install watcher`

### To Run

`watcher -config /path/to/config.json`

You may wish to redirect the output to a log file.

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

In this example when ever the file `/path/to/file/to/be/watched` is updated the file
`/path/to/file/to/be/uploaded` is uploaded to `http://url.to/upload/to` and `http://url.to/upload/to-stage`, 
also the file `/another/path/to/file/to/be/uploaded` is uploaded to `http://another.url.to/upload/to` and
`http://another.url.to/upload/to-stage`.
