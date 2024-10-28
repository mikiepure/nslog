# Non-structured logging

The nslog package provides non-structured logging, which is one of the implementation of "log/slog".

It is intended for logging to a console or text file as simple text line output.

Use the same interface as "log/slog". The most basic use and output is as follows;

```go
var logger = nslog.NewLogger(os.Stderr, nil)  // return slog.Logger object
logger.Info("log message", "key1", "val1")
// => 2024/10/31 11:22:33 INFO. log message key1=val1
//    ^^^^^^^^^^^^^^^^^^^ ^^^^^ ^^^^^^^^^^^ ^^^^^^^^^
//    time                level message     attrs-values
```

1st argument is text message and 2nd/3rd arguments is pair of key/value (Attrs-Values).
Attrs-Values can be omitted or more can be added.

```go
logger.Info("log message")
// => 2024/10/31 11:22:33 INFO. log message
logger.Info("log message", "key1", "val1", "key2", "val2")
// => 2024/10/31 11:22:33 INFO. log message key1=val1 key2=val2
```

Please see [Attrs and Values](#attrs-and-values) for more details.

The output and condition can be customized by LogHandlerOptions, the second argument of [nslog.NewLogger].
As an example,

```go
var logger = nslog.NewLogger(os.Stderr, &nslog.LogHandlerOptions{
    TimeLayout:     "2006/01/02 15:04:05.000",
    AddPID:         true,
    AddGoroutineID: true,
    AddSourceLevel: slog.LevelInfo,
    SourceFilePath: true,
})
logger.Info("log message")
// => 2024/10/31 11:22:33.444 ABCD 00000001 INFO. log message (d:/work/project/sample/main.go:19)
//    ^^^^^^^^^^^^^^^^^^^^^^^ ^^^^ ^^^^^^^^ ^^^^^ ^^^^^^^^^^^ ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
//    time                    pid  goroutineid    message     source
//                                          level
```

Please see [Options](#options) for more details.

## Options

The options listed in the following table can be configured to customize output of log message.

| Option         | Default Value         | Description |
| -------------- | --------------------- | ----------- |
| Level          | slog.LevelInfo        | Set level to output log message. By default, Error, Warn, and Info logs are output. |
| AddColor       | false                 | Add console color for level if it is true. |
| TimeLayout     | "2006/01/02 15:04:05" | Set own time layout for [Time.Format]. |
| AddPID         | false                 | Add PID as hex string if it is true. |
| AddGoroutineID | false                 | Add Goroutine ID as hex string if it is true. |
| AddSourceLevel | slog.LevelWarn        | Set level to output log source, which is the file and line number that called the function. By default, Error and Warn logs are output with source. |
| SourceFilePath | false                 | Use filepath for source if it is true. Use filename for source if it is false. |

These option can be overridden by environment variable.

| Option         | Environment Variable      | Available Value                             |
| -------------- | ------------------------- | ------------------------------------------- |
| Level          | GO_NSLOG_LEVEL            | "ERROR", "WARN", "INFO", or "DEBUG"         |
| AddColor       | GO_NSLOG_ADD_COLOR        | true: "TRUE" or "1" / false: "FALSE" or "0" |
| TimeLayout     | GO_NSLOG_TIME_LAYOUT      | Any string                                  |
| AddPID         | GO_NSLOG_ADD_PID          | true: "TRUE" or "1" / false: "FALSE" or "0" |
| AddGoroutineID | GO_NSLOG_ADD_GOROUTINEID  | true: "TRUE" or "1" / false: "FALSE" or "0" |
| AddSourceLevel | GO_NSLOG_ADD_SOURCE_LEVEL | "ERROR", "WARN", "INFO", or "DEBUG"         |
| SourceFilePath | GO_NSLOG_SOURCE_FILE_PATH | true: "TRUE" or "1" / false: "FALSE" or "0" |

## Groups

Groups can be added for logger.
As an examples,

```go
var logger = nslog.NewLogger(os.Stderr, nil).WithGroup("Main")
logger.Info("log message")
// => 2024/10/31 11:22:33 INFO. Main: log message
```

Attributes can also be added for logger.
As an examples,

```go
var logger = nslog.NewLogger(os.Stderr, nil).WithGroup("Main").With("id", 0)
logger.Info("log message")
// => 2024/10/31 11:22:33 INFO. Main[id=0]: log message
```

## Attrs and Values

Attrs and Values can be used for logging arguments. As an examples,

```go
var logger = nslog.NewLogger(os.Stderr, nil)
logger.Info("log message", "key", "val")
// => 2024/10/31 11:22:33 INFO. log message key=val
```
