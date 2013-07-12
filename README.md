Introduction
============

This package encodes arbitrary Go objects as XML property lists.

Type Mapping
============

| Go Type                       | XML Tag                 | Note                                  |
|-------------------------------|-------------------------|---------------------------------------|
| `bool`                        | `<true/>` or `</false>` |                                       |
| `int` and `uint` of any width | `<integer>`             | 					  |
| `float32` and `float64`       | `<real>`                | 					  |
| `string`                      | `<string>`              | 					  |
| `time.Time`                   | `<date>`                | RFC3339 formatted string 		  |
| `[]byte` and `[N]byte`        | `<data>`                | Base64 encoded string                 |
| `[]<T>` and `[N]<T>`          | `<array>`               |                                       |
| `struct` and `map`            | `<dict>`                | map keys will be converted to strings |

When encoding a `struct`, only the exported Fields are considered. By default the fields name is
used as the dictionary key. A custom name may be chosen, by setting the `plist`-key of the fields
tag. Setting the tag to `"-"` causes the field to be ignored.


Installation
============

This package can be installed using:

    go get github.com/wkz/plist

Usage
=====

See examples/