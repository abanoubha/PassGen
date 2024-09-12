package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/exp/rand"
)

var (
	version   bool
	help      bool
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

	rootCmd.Flags().BoolVarP(&help, "help", "h", false, "show help screen")

	rootCmd.Flags().StringVarP(&count, "count", "c", "", "calculate the frequency of occurrence of each character in a plain text file")

	rootCmd.Flags().IntVarP(&password, "password", "p", 0, "generate passwords according to the character frequency specified")
	rootCmd.Flags().StringVarP(&charFreqs, "charfreq", "f", "", "specify the character frequency table to be used")

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		if version {
			printHelpScreen()
		} else if help {
			printHelpScreen()
		} else if count != "" {
			countCharFreq(count)
		} else if password > 0 && charFreqs != "" {
			generatePasswords(password, charFreqs)
		} else {
			fmt.Println(`You did not used the correct arguments.`)
			printHelpScreen()
		}
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func printHelpScreen() {
	fmt.Println(`
pwdgen v0.1.0 beta
	
Software Developer  : Abanoub Hanna
Source Code         : https://github.com/abanoubha/PassGen
X platform          : https://x.com/@AbanoubHA
Developer's Website : https://AbanoubHanna.com

Examples:
  pwdgen -v # show the app version
  pwdgen -c textfile.txt # show the count of each character occurrences
  pwdgen -p 8 -f charfreq.txt # generate passwords
  pwdgen -h # show (this) help screen`)
}

type CharFreq struct {
	Char string
	Freq int
}

type KV struct {
	Key string
	Val int
}

func countCharFreq(filename string) {
	fmt.Println("counting frequency of character occurrences in ", filename, " ...")

	counts := countChars(filename)

	kvPairs := make([]KV, 0, len(counts))

	for k, v := range counts {
		kvPairs = append(kvPairs, KV{k, v})
	}

	sort.Slice(kvPairs, func(i, j int) bool {
		return kvPairs[i].Val > kvPairs[j].Val
	})

	for char, count := range kvPairs {
		fmt.Printf("%v: %v\n", char, count)
	}
}

func countChars(filename string) map[string]int {
	charCounts := make(map[string]int)
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		for _, char := range scanner.Text() {
			charCounts[string(char)]++
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return charCounts
}

func generatePasswords(passwordLength int, charFreqsFile string) {
	fmt.Println("generating passwords with length of ", passwordLength, " respecting the character freqency specified in ", charFreqsFile, " ...")

	charFreqs, err := readCharFreqs(charFreqsFile)
	if err != nil {
		panic("can not read the character frequency file specified.")
	}

	cumFreqs := make([]int, len(charFreqs))
	cumFreqs[0] = charFreqs[0].Freq
	for i := 1; i < len(charFreqs); i++ {
		cumFreqs[i] = cumFreqs[i-1] + charFreqs[i].Freq
	}

	for i := 1; i < 1_000_000_000; i++ {
		fmt.Println(generateStrings(charFreqs, cumFreqs, passwordLength))
	}
}

func generateStrings(charFreqs []CharFreq, cumFreqs []int, passLength int) string {
	rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	// Calculate total frequency
	totalFreq := 0
	for _, freq := range charFreqs {
		totalFreq += freq.Freq
	}

	// Generate strings
	// strings := make([]string, numStrings)
	// for i := 0; i < numStrings; i++ {
	var sb bytes.Buffer
	for j := 0; j < passLength; j++ {
		// Generate a random number between 0 and totalFreq-1
		randNum := rand.Intn(totalFreq)

		// Find the corresponding character based on the cumulative frequency distribution
		index := 0
		for ; index < len(cumFreqs); index++ {
			if randNum < cumFreqs[index] {
				break
			}
		}
		sb.WriteString(charFreqs[index].Char)
	}
	// strings[i] = sb.String()
	// }

	// return strings
	return sb.String()
}

func readCharFreqs(charFreqsFile string) ([]CharFreq, error) {
	csvFile, err := os.Open(charFreqsFile)
	if err != nil {
		return []CharFreq{}, fmt.Errorf("error opening file: %v", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	// _, err = reader.Read()
	// if err != nil {
	// 	fmt.Println("Error reading header:", err)
	// 	return
	// }

	charFreqs := []CharFreq{}

	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break // End of file
			}
			return []CharFreq{}, fmt.Errorf("error reading CSV file records : %v", err)
		}

		freq, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, err
		}

		charFreqs = append(charFreqs, CharFreq{record[0], freq})
	}

	return charFreqs, nil
}
