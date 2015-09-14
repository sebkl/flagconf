# flagconf
Package flagconf wraps the standard golang package [flag](https://golang.org/pkg/flag/). It intercepts it and preloads configuration values from
local files (home-directory, /etc) and environment variables in order to retrieve default values. The resulting
 order of flag loading is:

1. in-code values (defaults provided in [flagset](https://golang.org/pkg/flag/#FlagSet), which is usually defined in the init function)
2. local file list *(default is "~/.falgconf/<PREFIX>.yaml")*
3. environment variables
4. command line arguments

## Usage:
```go
package main

import (
	"github.com/sebkl/flagconf"
	"flag"
)

var val *string
func init() {
	val = flag.String("val","defaultvalue","This is my application flag."
}

func main() {
	// The file "~/.flagconf/MYAPP.yml" and
	// all environemnt variables starting with "MYAPP"
	// are being evaluated for a flag value named "val"
	// If found the default value "defaultvalue" is overwritten.
	flagconf.Parse("MYAPP")

	/* ... code ... */
}
```

## Features

- Json encoded configuration files instead of yaml.
- Provide an ordered list of files to look for flag values by the [FileList](http://godoc.org/github.com/sebkl/flagconf#FileList) function.

### Additional flags added by flagconf
| flag | Description |
| --- | --- |
| ```flagconfConfirm``` | Show all flag configs before execution. |
| ```flagconfFile```| Explicitly define a file to read flag configs from. |

## Documentation
Please find documentation at [godoc.org/github.com/sebkl/flagconf](http://godoc.org/github.com/sebkl/flagconf).

## Contribution
Pleas feel free to send in Pull request. I want to make this a more comprehensive library for dealing with command line flags.
