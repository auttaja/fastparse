# fastparse
A fast argument parser for Go designed for situations where you know the maximum amount of memory which a message can use (for example, a Discord or Twitter message). This library is extremely quick since it tries to avoid allocating memory on the fly, instead opting for pre-allocating the memory which the application will need. To use this library, when you initialise your application/library which uses this you will want to create a parser manager:
```go
// Creates the parser manager. This will be used to create child parsers.
// The first parameter is the maximum length of a message. In this example, we are using Discord's limit of 2000 characters.
// The second parameter is the number of pre-allocated pads which will be used.
// Each pre-allocated pad will use roughly the length of a message, but will reduce the chance that a memory allocation needs to be made because all pads are being used.
m := fastparse.NewParserManager(2000, 100)
```

From here, we can simply create a parser from a thread when we need it. The reader must be a `io.ReadSeeker`:
```go
parser := m.Parser(reader)
```

**Please note that you should mark your parser as done when complete. This is important since if it's pre-allocated, it will not be re-added to the pool otherwise. You can do this with the `Done` function:**
```go
defer parser.Done()
```

The parser has the following other functions which you can use:
- `Remainder() (string, error)`: Gets the remainder of a reader.
- `GetNextArg() *Argument`: GetNextArg is used to get the next argument. If there are no additional arguments, the pointer will be nil. The Argument struct contains the following:
    - `Text`: The text from the argument as a string.
    - `Rewind() error`: Rewind is used to rewind the reader to before an argument was read. This is useful for some argument verification situations. Note that if you want to rewind multiple arguments, you need to run it on every argument you parsed AFTER the one you wish to rewind in order of last to first.
