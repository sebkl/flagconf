package flagconf_test

import (
	flagconf "."
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var jsonflag *string

func init() { /* initialize to flagset as usual ... */
	// Usual configuration of the flagset using the flag package
	jsonflag = flag.String("jsonflag", "defaultvalue", "Just a test flag")
}

func Example_json() {
	// -- set an exemplary variable in filesystem
	ioutil.WriteFile(TMPFILE+".json", []byte(`{ "jsonflag": "jsonvalue"}`), 0644)
	flagconf.FileList([]string{TMPFILE + ".json"})
	defer os.Remove(TMPFILE + ".json") // just make sure to remove temporary file
	// ---
	// The default used file is : "~/.flagconf/<PREFIX>.yml"

	//Actual code in main() function:
	flagconf.Parse(PREFIX) // Compared to the flag package, the prefix is the only additional information passed to the parse function.

	fmt.Printf("%s", *jsonflag)

	//Output:
	//jsonvalue
}
