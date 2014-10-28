# The Watcher

### Requirements

- Go (`brew install go` for macos with brew)
- Go cross compilation support (see https://github.com/davecheney/golang-crosscompile)

### To Build

- `GOOS=windows GOARCH=amd64 go install watcher`
- `GOOS=windows GOARCH=386 go install watcher`

### To Run

`watcher -config /path/to/config.json`

You may wish to redirect the output to a log file.

### Example Config File

```json
{
  "log_file": "/path/to/log/file.log",
  "rollbar_token": "a2cf259b377...",
  "production": true,
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
`http://another.url.to/upload/to-stage`. Additionally, any errors will be logged to rollbar with the supplied 
rollbar token.  Other info/error information will the logged to the specified log file.
