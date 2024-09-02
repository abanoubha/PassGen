package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version   bool
	count     string
	charFreqs string
	password  int
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "passgen",
		Short: "password generator based on character frequency",
		Long:  "password generator based on character frequency",
		Example: `passgen -c rockyou.txt # get counts of all character in a plain text file
passgen -p 8 # generate passwords (continuously) with 8 characters length
passgen -v # show version number and release info`,
	}

	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "show version number and release info")

	rootCmd.Flags().StringVarP(&count, "count", "c", "", "calculate the frequency of occurrence of each character in a plain text file")

	rootCmd.Flags().IntVarP(&password, "password", "p", 0, "generate passwords according to the character frequency specified")
	rootCmd.Flags().StringVarP(&charFreqs, "charfreq", "f", "", "specify the character frequency table to be used")

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		if version {
			fmt.Println(`
passgen v0.1.0 beta

Software Developer  : Abanoub Hanna
Source Code         : https://github.com/abanoubha/PassGen
X platform          : https://x.com/@AbanoubHA
Developer's Website : https://AbanoubHanna.com`)
		} else if count != "" {
			countCharFreq(count)
		} else if password > 0 && charFreqs != "" {
			generatePasswords(password, charFreqs)
		}
	}
}

func countCharFreq(filename string) {
	fmt.Println("counting frequency of character occurrences in ", filename, " ...")
}

func generatePasswords(passwordLength int, charFreqsFile string) {
	fmt.Println("generating passwords with length of ", passwordLength, " respecting the character freqency specified in ", charFreqsFile, " ...")
}
