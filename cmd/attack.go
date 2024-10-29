/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/jacobmaniscalco/goprobe/internal/attack"
	ssh "github.com/jacobmaniscalco/goprobe/internal/attack/modules/ssh"
)

var attackOptions attack.AttackOptions 

// attackCmd represents the attack command
var attackCmd = &cobra.Command{
	Use:   "attack",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := ssh.BruteForceSSH(attackOptions)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(attackCmd)

	attackCmd.Flags().StringVarP(&attackOptions.Host, "target", "t", "",
	"Specify the target IP address or range of IP addresses to attack." +
	"This can be a single IP, a subnet, or a list of IPs.")

	attackCmd.Flags().StringVarP(&attackOptions.Port, "port", "p","22",
	"Specify the port to attack") 

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// attackCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// attackCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
