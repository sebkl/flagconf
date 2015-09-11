package flagconf_test

import (
	flagconf "."
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var fsflag *string

func init() { /* initialize to flagset as usual ... */
	// Usual configuration of the flagset using the flag package
	fsflag = flag.String("fsflag", "defaultvalue", "Just a test flag")
}

func Example_fs() {
	// --- set an exemplary variable in filesystem
	ioutil.WriteFile(TMPFILE+".yml", []byte("fsflag: myvalue"), 0644)
	flagconf.FileList([]string{TMPFILE + ".yml"})
	defer os.Remove(TMPFILE + ".yml") // just make sure to remove temporary file
	// ---
	// The default used file is : "~/.flagconf/<PREFIX>.yml"

	//Actual code in main() function:
	flagconf.Parse(PREFIX)
	// the only additional information compared to is the application specific PREFIX t

	fmt.Printf("%s", *fsflag)

	//Output:
	//myvalue
}
