package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var echoTimes int

var cmdRemoveImagesDangling = &cobra.Command{
	Use:   "remove-images-dangling",
	Short: "Remove images dangling",
	Long:  `remove all images "dangling=true"`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Print: " + strings.Join(args, " "))
	},
}

var cmdContainerStatusExited = &cobra.Command{
	Use:   "remove-containers-exited [string to remove]",
	Short: "Remove containers exited",
	Long:  `remove all containers "status=exited"`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Print: " + strings.Join(args, " "))
	},
}

//var cmdListImages = &cobra.Command{
//	Use:   "list-images [string to remove]",
//	Short: "List images",
//	Long:  `list all images`,
//	Args:  cobra.MinimumNArgs(1),
//	Run: func(cmd *cobra.Command, args []string) {
//		fmt.Println("Print: " + strings.Join(args, " "))
//	},
//}

var cmdEcho = &cobra.Command{
	Use:   "echo [string to echo]",
	Short: "Echo anything to the screen",
	Long: `echo is for echoing anything back.
Echo works a lot like print, except it has a child command.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Echo: " + strings.Join(args, " "))
	},
}

var cmdTimes = &cobra.Command{
	Use:   "times [string to echo]",
	Short: "Echo anything to the screen more times",
	Long: `echo things multiple times back to the user by providing
a count and a string.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for i := 0; i < echoTimes; i++ {
			fmt.Println("Echo: " + strings.Join(args, " "))
		}
	},
}

func AddCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}

var rootCmd = &cobra.Command{Use: "gdocker"}

func Execute() {
	cmdTimes.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo the input")

	rootCmd.AddCommand(cmdEcho, cmdRemoveImagesDangling, cmdContainerStatusExited)
	cmdEcho.AddCommand(cmdTimes)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
