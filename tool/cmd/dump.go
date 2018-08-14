// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sapk/genesys/api/client"
	"github.com/sapk/genesys/api/object"
	"github.com/sapk/genesys/tool/check"
)

//TODO add flag for username and pass
//TODO add short option that only dump json doesn't format data
//TODO add from dump (from short flag
//TODO flag for no cleanup before export

// listCmd represents the list command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Connect to a GAX server to dump his state",
	Args: func(cmd *cobra.Command, args []string) error {
		logrus.Debug("Checking args for list cmd: ", args)
		if len(args) > 1 {
			return fmt.Errorf("requires at least one GAX server")
		}
		for _, arg := range args {
			if !check.IsValidClientArg(arg) {
				return fmt.Errorf("invalid argument specified (ex: gax_host:8080): %s", arg)
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		for _, gax := range args {
			if !strings.Contains(gax, ":") {
				//By default use port 8080
				gax += "8080"
			}

			///Login
			c := client.NewClient(gax)
			user, err := c.Login("default", "password")
			if err != nil {
				logrus.Panicf("Login failed : %v", err)
			}
			logrus.Println(user)

			//Cleanup
			err = cleanAll()
			if err != nil {
				logrus.Panicf("Clean up failed : %v", err)
			}

			//Get DATA
			//Hosts
			hosts, err := c.ListHost()
			if err != nil {
				logrus.Panicf("ListHost failed : %v", err)
			}
			sort.Sort(object.CfgHostList(hosts)) //order data by name
			err = dumpToFile("./Hosts.json", hosts)
			if err != nil {
				logrus.Panicf("Dump failed : %v", err)
			}
			//Applications
			apps, err := c.ListApplication()
			if err != nil {
				logrus.Panicf("ListApplication failed : %v", err)
			}
			sort.Sort(object.CfgApplicationList(apps)) //order data by name
			err = dumpToFile("./Applications.json", apps)
			if err != nil {
				logrus.Panicf("Dump failed : %v", err)
			}

			//TODO format data
			err = os.Mkdir("Host", 0755)
			if err != nil {
				logrus.Panicf("Folder creation failed : %v", err)
			}
			err = os.Mkdir("Applications", 0755)
			if err != nil {
				logrus.Panicf("Folder creation failed : %v", err)
			}
			for _, host := range hosts {
				logrus.Infof("Host: %s (%s)", host.Name, host.DBID)
				err = writeToFile("Host/"+host.Name+".md", fmt.Sprintf("Host: %s (%s)", host.Name, host.DBID))
				if err != nil {
					logrus.Panicf("File creation failed : %v", err)
				}
			}
			for _, app := range apps {
				logrus.Infof("App: %s (%s) @ %s", app.Name, app.DBID, app.WorkDirectory)
				err = writeToFile("Applications/"+app.Name+".md", fmt.Sprintf("App: %s (%s) @ %s", app.Name, app.DBID, app.WorkDirectory))
				if err != nil {
					logrus.Panicf("File creation failed : %v", err)
				}
			}
		}
	},
}

func writeToFile(file, data string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(data)
	if err != nil {
		return err
	}
	f.Sync()
	return nil
}

func cleanAll() error {
	err := os.RemoveAll("./Hosts.json")
	if err != nil {
		return err
	}
	err = os.RemoveAll("./Applications.json")
	if err != nil {
		return err
	}
	err = os.RemoveAll("./Hosts")
	if err != nil {
		return err
	}
	return os.RemoveAll("./Applications")
}

func dumpToFile(file string, data interface{}) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, json, 0644)
}

func init() {
	RootCmd.AddCommand(dumpCmd)
}
