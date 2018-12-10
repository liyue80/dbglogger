# dbglogger

Simple RESTful server which write the incoming message to disk file or print to console window.

## Sample

curl -X POST -H "Content-Type: application/json" -d @m1.json http://127.0.0.1:27109/dbgloggers

## Configure

`PrintConsole`

Boolean. Show incoming messages on stdout.

`ConsoleSeverity`

Integer. Print to stdout if incoming message severity is equal or lower than this threshold, except this threshold value is zero, which means printing all incoming message without check therir severity.

`PrintFile`

Boolean. Write incoming message to disk file.

`FileName`

String. File path name.

`FileSeverity`

Integer. Refer to _ConsoleSeverity_.
