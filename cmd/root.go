package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/websocket"

	"github.com/hypebeast/go-osc/osc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "minion-sonicpi",
	Short: "A minion to drive the Sonic-Pi app",
	Run: func(cmd *cobra.Command, args []string) {

		webstrate := viper.Get("webstrate")

		if webstrate == "" {
			cmd.Usage()
			return
		}

		origin := "https://emet.cc.au.dk/"
		url := fmt.Sprintf("wss://emet.cc.au.dk/minion/v1/connect/%s?type=minion-sonic", webstrate)
		ws, err := websocket.Dial(url, "", origin)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := ws.Write([]byte("drop the beat")); err != nil {
			log.Fatal(err)
		}

		client := osc.NewClient("localhost", 4557)

		for {
			var msg = make([]byte, 512)
			if _, err = ws.Read(msg); err != nil {
				log.Fatal(err)
				break
			}
			smsg := string(msg)
			if strings.HasPrefix(smsg, "stop") {
				stopmsg := osc.NewMessage("/stop-all-jobs")
				stopmsg.Append(int32(111))
				client.Send(stopmsg)
				fmt.Println("stopped")
			} else {
				runmsg := osc.NewMessage("/run-code")
				runmsg.Append(int32(111))
				runmsg.Append(smsg)
				client.Send(runmsg)
				fmt.Println("playing: " + smsg)
			}
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.minion-sonicpi.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().StringP("webstrate", "w", "", "Id of the webstrate you want me to connect to")
	RootCmd.MarkFlagRequired("webstrate")
	viper.BindPFlags(RootCmd.Flags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".minion-sonicpi") // name of config file (without extension)
	viper.AddConfigPath("$HOME")           // adding home directory as first search path
	viper.AutomaticEnv()                   // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
