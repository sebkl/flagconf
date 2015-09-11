# flagconf
Package flagconf wraps the standard golang package [flag](https://golang.org/pkg/flag/)flag. It intercepts it and preloads configuration values from
local files (home-directory, /etc) and environment variables in order to retrieve default values. The resulting
 order of flag loading is:

1. in-code values (defaults provided in [flagset](https://golang.org/pkg/flag/#FlagSet), which is usually defined in the init function)
2. local file list *(default is just ~/.falgconf/<PREFIX>.yaml)*
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
	val flag.String("val","defaultvalue","This is my application flag."
}

func main() {
	//The file "~/.flagconf/MYAPP.yml" is evaluated for the flag named val and if existing will overwrite the "defaultvalue"
	flag.Parse("MYAPP")

	//	
}
```

## Features

- Json encoded configuration files instead of yaml.
- Provide an orderd list of files to look for flag values in using the [FileList](http://godoc.org/sebkl/flagconf#FileList) function.


## Documentation
Please find documentation at [godoc.org/sebkl/flagconf](http://godoc.org/sebkl/flagconf).

# Contribution
Pleas feel free to send in Pull request. I want to make this a more comprehensive library for dealing with command line flags.
