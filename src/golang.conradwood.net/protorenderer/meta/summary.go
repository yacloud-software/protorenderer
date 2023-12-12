package meta

import (
	"fmt"
	"strings"
)

func printPackage(pkg *Package) {
	for _, s := range pkg.Services {
		fmt.Printf("Service: %s\n", s.Name)
		printPrefix("         ", "%s\n", s.Comment)
		for _, r := range s.RPCs {
			fmt.Printf("   RPC: %s\n", r.Name)
			printPrefix("       ", r.Comment)
		}
	}
	for _, m := range pkg.Messages {
		fmt.Printf("Message: %s\n", m.Name)
		printPrefix("          ", "%s\n", m.Comment)
		for _, f := range m.Fields {
			if len(strings.Split(f.Comment, "\n")) < 2 {
				fmt.Printf("  Field: %s %s // %s\n", f.TypeName(), f.Name, f.Comment)
			} else {
				fmt.Printf("  Field: %s %s\n", f.TypeName(), f.Name)
				printPrefix("          ", "%s\n", f.Comment)
			}
		}
	}
}

func printPrefix(prefix string, txt string, args ...interface{}) {
	s := fmt.Sprintf(txt, args...)
	s = strings.Trim(s, "\n")
	for _, line := range strings.Split(s, "\n") {
		fmt.Printf("%s%s\n", prefix, line)
	}
}




























































































