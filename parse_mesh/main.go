// Command parse_mesh is an ASCII MeSH format parser. (It's actually UTF8 but for
// some reason they call it ASCII and rename it .bin instead of .txt)
//
// It takes three command-line arguments:
//   1) the input ascii-mesh file
//   2) the nodes output file (3-column tab-delimited)
//   3) the edges output file (3-column tab-delimited, for ancestor relationships)
//
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "USAGE: %s asciimesh.bin nodes.txt edges.txt\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  where nodes.txt will receive the node listing\n")
		fmt.Fprintf(os.Stderr, "    and edges.txt will receive the self edges listing\n")
		return
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	nfile, err := os.OpenFile(os.Args[2], os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	efile, err := os.OpenFile(os.Args[3], os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer nfile.Close()
	defer efile.Close()

	id := ""   // MeSH UI
	name := "" // MeSH Heading "MH"
	def := ""  // MeSH Subject "MS"
	lastkey := ""
	rectype := ""
	mnums := make([]string, 0, 10)

	// map from MN to UI
	treenums := make(map[string]string)

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}
		if line == "*NEWRECORD" {
			if id != "" {
				fmt.Fprintf(nfile, "%s\t%s\t%s\n", id, name, def)
				for _, mn := range mnums {
					treenums[mn] = id
				}
			}
			id = ""
			name = ""
			def = ""
			lastkey = ""
			rectype = ""
			mnums = mnums[:0]
			continue
		}

		row := strings.SplitN(line, " = ", 2)
		if len(row) != 2 {
			// this is wonky, but it seems that some MS lines wrap in d2017...
			if lastkey == "MS" {
				def += " " + strings.TrimSpace(strings.Replace(line, "\t", " ", -1))
			}
			continue
		}
		switch row[0] {
		case "RECTYPE":
			rectype = strings.TrimSpace(row[1])

		case "UI":
			id = strings.TrimSpace(row[1])

		case "NM":
			if rectype == "C" {
				name = strings.TrimSpace(row[1])
				// chemicals don't seem to have a description, so we duplicate name
				// should we just leave it blank instead?
				def = name
			}

		case "MH":
			if rectype == "D" {
				name = strings.TrimSpace(row[1])
			}

		case "SH":
			if rectype == "Q" {
				name = strings.TrimSpace(row[1])
			}

		case "MS": // "Q" and "D" rectype
			def = strings.TrimSpace(strings.Replace(row[1], "\t", " ", -1))

		case "MN":
			mnums = append(mnums, strings.TrimSpace(row[1]))
		}
		lastkey = row[0]
	}

	if id != "" {
		fmt.Fprintf(nfile, "%s\t%s\t%s\n", id, name, def)
		for _, mn := range mnums {
			treenums[mn] = id
		}
	}

	//////////////

	for mn, ui := range treenums {
		idx := strings.LastIndex(mn, ".")
		if idx == -1 {
			// top-level term, no ancestors
			continue
		}
		anc := mn[:idx]
		if other, ok := treenums[anc]; ok {
			fmt.Fprintf(efile, "%s\t%s\t%s\n", ui, other, "has_broader_term")
		} else {
			panic("unknown mesh number " + anc)
		}
	}
}
