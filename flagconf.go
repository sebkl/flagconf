package flagconf

// Package flagconf wraps the standard golang package flag. It intercepts it and preloads configuration values from
// local files (home-directory, /etc) and environment variables in order to retrieve default values. The resulting
// order of flag loading is:
//
//  1) in-code values (defaults provided in flagset, whcih is usually defined in the init function)
//  2) local file list (default is just ~/.falgconf/<PREFIX>.yaml
//  3) environment variables
//  4) command line arguments
//
// TODO: implement option that dumps loaded options and asks the user whether to really use them.

import (
	"encoding/json"
	"flag"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
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

var filelist []string

//FileList sets or gets the list of used files in actual order.
func FileList(fl ...[]string) []string {
	if len(fl) > 0 {
		filelist = fl[0]
	}
	return filelist
}

//Parse is a replacement for the flag.Parse function that intercepts it and reads the configuration
// from a defined set of files and environment variables.
func Parse(prefix string, ofs ...*flag.FlagSet) {
	if len(ofs) == 0 {
		ofs = []*flag.FlagSet{flag.CommandLine}
	}

	for _, fs := range ofs {
		if filelist == nil { // if not filelist is set by user, use a default one
			filelist = []string{"~/.flagconf/" + prefix + ".yml"}
		}

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
}