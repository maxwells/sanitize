package main

import (
	"os"
	"fmt"
	"github.com/maxwells/sanitizer/sanitize"
)

func main() {
	whitelist, _ := sanitize.WhitelistFromFile("./whitelist.json")

	f, _ := os.Open("./example.html")
	sanitized, _ := whitelist.SanitizeRemove(f)
	fmt.Printf("%s\n",sanitized)
}