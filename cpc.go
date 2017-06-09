package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/ghodss/yaml"
	"github.com/urfave/cli"
	"github.com/davecgh/go-spew/spew"
)

type appFlags struct {
	BoolFlags   []cli.BoolFlag
	StringFlags []cli.StringFlag
	IntFlags    []cli.IntFlag
}

var flags *appFlags

func init() {
	// Find YAML file in project
	yamlFile, err := ioutil.ReadFile("cpc.yaml")
	if err != nil {
		fmt.Println("ERR: ", err)
	}
	err = yaml.Unmarshal(yamlFile, &flags)
	if err != nil {
		fmt.Println("ERR: ", err)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "cpc"
	app.Version = "0.1.0"
	app.HideHelp = true
	app.HideVersion = true
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		{
			Name:  "Luke Hobbs",
			Email: "lukeehobbs@gmail.com",
		},
	}
	app.Usage = "Control CI pipeline using commit messages."
	app.UsageText = "commit -m \"<your commit message> cpc [options] [arguments]\""
	app.Action = func(c *cli.Context) error {
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)
		for _, b := range flags.BoolFlags {
			bs := strconv.FormatBool(c.IsSet(b.Name))
			fmt.Fprintf(w, "%s:\t%s\n", b.Name, bs)
			os.Setenv(strings.ToUpper(b.Name), bs)
		}
		for _, s := range flags.StringFlags {
			ss := c.String(s.Name)
			fmt.Fprintf(w, "%s:\t%s\n", s.Name, ss)
			os.Setenv(strings.ToUpper(s.Name), ss)
		}
		for _, i := range flags.IntFlags {
			is := strconv.Itoa(c.Int(i.Name))
			fmt.Fprintf(w, "%s:\t\t\t%s\n", i.Name, is)
			os.Setenv(strings.ToUpper(i.Name), is)
		}
		w.Flush()
		// TODO: Create yaml/json file containing these variables for the pipeline to reference
		return nil
	}
	app.Flags = []cli.Flag{
	// Flags specific to cpc will go here
	}

	for _, b := range flags.BoolFlags {
		app.Flags = append(app.Flags, b)
	}
	for _, s := range flags.StringFlags {
		app.Flags = append(app.Flags, s)
	}
	for _, i := range flags.IntFlags {
		app.Flags = append(app.Flags, i)
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	for i, s := range os.Args {
		if s == "cpc" {
			os.Args = os.Args[i:]
		}
	}
	app.Run(os.Args)
}
