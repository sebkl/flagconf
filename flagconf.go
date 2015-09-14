// Package flagconf wraps the standard golang package flag. It intercepts it and preloads configuration values from
// local files (home-directory, /etc) and environment variables in order to retrieve default values. The resulting
// order of flag loading is:
//
//  1) in-code values (defaults provided in the flagset)
//  2) local file list (default is just [ "~/.flagconf/<PREFIX>.yaml" ]
//  3) environment variables
//  4) command line arguments
//
// TODO: implement option that dumps loaded options and asks the user whether to really use them.
package flagconf

import (
	"encoding/json"
	"flag"
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
)

func setFromFileByPrefix(prefix, fn string, ret map[string]string) {
	p := path.Ext(fn)
	_, err := os.Stat(fn)
	if err != nil {
		return
	}

	file, err := os.Open(fn)
	if err != nil {
		log.Printf("Could not open flagconf file '%s': %s", fn, err)
		return
	}

	by, err := ioutil.ReadAll(file)
	if err != nil {
		goto out
	}

	switch p {
	case ".json":
		err = json.Unmarshal(by, &ret)
	default: // assume yaml
		err = yaml.Unmarshal(by, &ret)
	}

out:
	if err != nil {
		log.Printf("flagconf (prefix: %s, suffix: %s): %s", prefix, p, err)
	}
}

func setFromEnvByPrefix(prefix string, ret map[string]string) {
	env := os.Environ()
	for _, v := range env {
		vals := strings.Split(v, "=")
		if len(vals) < 2 {
			continue
		}

		key := vals[0]
		val := vals[1]

		if strings.HasPrefix(key, prefix) {
			fkey := strings.TrimLeft(key, prefix)
			ret[fkey] = val
		}
	}
}

func confirmFlagConfiguration(fss ...*flag.FlagSet) bool {
	for _, fs := range fss {
		fs.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(os.Stderr, "\t%s = %s\n", f.Name, f.Value.String())
		})
	}
	fmt.Fprintf(os.Stderr, "\nProceed with the above flag configuration ? (yes/[no]) ")
	var resp string
	_, err := fmt.Scanln(&resp)
	if err != nil {
		return false
	}

	if ok, err := regexp.MatchString("^[yY]([eE][sS])?$", resp); ok && err == nil {
		return true
	}
	return false
}

var filelist []string
var internalConfig struct {
	ConfigFile          string
	RequireConfirmation bool
}

//FileList sets or gets the list of used files in actual order.
func FileList(fl ...[]string) []string {
	if len(fl) > 0 {
		filelist = fl[0]
	}
	return filelist
}

func setupInternalFlags(fs *flag.FlagSet) {
	if fs.Lookup("flagconfFile") == nil {
		fs.StringVar(&internalConfig.ConfigFile, "flagconfFile", "", "Specify flagconf configuration")
	}

	if fs.Lookup("flagconfConfirm") == nil {
		fs.BoolVar(&internalConfig.RequireConfirmation, "flagconfConfirm", false, "Require confirmation from user ro use flag configuration")
	}
}

//Parse is a replacement for the flag.Parse function that intercepts it and reads the configuration
// from a defined set of files and environment variables.
func Parse(prefix string, ofs ...*flag.FlagSet) {
	if len(ofs) == 0 {
		ofs = []*flag.FlagSet{flag.CommandLine}
	}

	if filelist == nil { // if not filelist is set by user, use a default one
		home := os.Getenv("HOME")
		filelist = []string{home + "/.flagconf/" + prefix + ".yml"}
	}

	setupInternalFlags(ofs[0])

	for _, fs := range ofs {

		settings := make(map[string]string)

		//walk through levels
		// 4th level is manaaged by flags (defaults in flagset)
		// 3rd level: filelist
		for _, fn := range filelist {
			setFromFileByPrefix(prefix, fn, settings)
		}
		// 2nd level: environment variables
		setFromEnvByPrefix(prefix, settings)

		//1st and lowest (4th) level is managed by flags
		fs.VisitAll(func(f *flag.Flag) {
			if val, ok := settings[f.Name]; ok {
				f.Value.Set(val)
			}
		})
	}

	flag.Parse()

	if internalConfig.RequireConfirmation {
		if !confirmFlagConfiguration(ofs...) {
			log.Fatal("Aborted. Flag configuration not accepted by user.")
		}
	}
}
