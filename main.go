package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		maskedLine := maskPublicIPs(line)
		fmt.Println(maskedLine)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
	}
}

func maskPublicIPs(input string) string {
	ipv4Regex := `\b(?:\d{1,3}\.){3}\d{1,3}\b`
	re := regexp.MustCompile(ipv4Regex)
	matches := re.FindAllString(input, -1)

	for _, match := range matches {
		if isPublicIPv4(match) {
			input = strings.Replace(input, match, "XXX.XXX.XXX.XXX", -1)
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
