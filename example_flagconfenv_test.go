package flagconf_test

import (
	"github.com/sebkl/flagconf"
	"flag"
	"fmt"
	"os"
)

var envflag *string

func init() { /* initialize to flagset as usual ... */
	envflag = flag.String("envflag", "defaultvalue", "Just a test flag.")

	// set an exemplary environment variable that matches the prefix pattern
	os.Setenv(PREFIX+"envflag", "testvalue") //set env variable
}

func Example_basic() {
	//Actual code in main() function:
	flagconf.Parse(PREFIX)
	// the only additional information compared to is the application specific PREFIX t

	fmt.Printf("%s", *envflag)

	//Output:
	//testvalue
}
