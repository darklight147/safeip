package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var maskDNS bool
var customMask string
var customDnsRegex string

var rootCmd = &cobra.Command{
	Use:   "safeip",
	Short: "Mask public IPs and DNS-like entries",
	Run: func(cmd *cobra.Command, args []string) {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			maskedLine := maskPublicIPs(line)
			fmt.Println(maskedLine)
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
		}
	},
	Long: `SafeIP is a command-line tool written in Go that helps you mask public IPv4 addresses and DNS-like entries from your text input. It's useful for redacting sensitive information from logs or other textual data to avoid people eyeballing your data 🕵️`,
}

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates bash completion scripts",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(shellCompletionFunction(args[0]))
	},
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"bash", "zsh", "fish"},
}

func shellCompletionFunction(shell string) string {
	// Define the completion script
	switch shell {
	case "bash":
		rootCmd.GenBashCompletion(os.Stdout)
		break
	case "zsh":
		rootCmd.GenZshCompletion(os.Stdout)
		break
	case "fish":
		rootCmd.GenFishCompletion(os.Stdout, true)
		break

	}

	return ""

}

func init() {
	rootCmd.Flags().BoolVar(&maskDNS, "mask-dns", false, "Mask DNS-like entries")
	rootCmd.Flags().StringVar(&customMask, "mask", "XXX.XXX.XXX.XXX", "Custom mask to use")
	rootCmd.Flags().StringVar(&customDnsRegex, "dns-regex", `(\b(?:[a-zA-Z0-9-]+\.){2,}[a-zA-Z]{2,}\b)`, "Custom regex to use for DNS-like entries")
	rootCmd.AddCommand(completionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func maskPublicIPs(input string) string {
	ipv4Regex := `\b(?:\d{1,3}\.){3}\d{1,3}\b`
	ipv6Regex := `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`
	re := regexp.MustCompile(ipv4Regex)
	re6 := regexp.MustCompile(ipv6Regex)
	matches := re.FindAllString(input, -1)
	ipv6Matches := re6.FindAllString(input, -1)

	for _, match := range matches {
		if isPublicIPv4(match) {
			input = strings.Replace(input, match, customMask, -1)
		}
	}

	for _, match := range ipv6Matches {
		input = strings.Replace(input, match, customMask, -1)
	}

	if maskDNS {
		dnsRegex := customDnsRegex
		dnsRe := regexp.MustCompile(dnsRegex)
		dnsMatches := dnsRe.FindAllString(input, -1)
		for _, match := range dnsMatches {
			if !strings.Contains(match, ".internal") {
				input = strings.Replace(input, match, customMask, -1)
			}
		}
	}

	return input
}

func isPublicIPv4(ip string) bool {
	parts := strings.Split(ip, ".")
	firstPart, _ := strconv.Atoi(parts[0])
	secondPart, _ := strconv.Atoi(parts[1])
	if firstPart == 10 || firstPart == 172 && (16 <= secondPart && secondPart <= 31) || firstPart == 192 && parts[1] == "168" {
		return false
	}
	return true
}
