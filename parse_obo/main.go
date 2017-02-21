// Command parse_obo is a very simple OBO format parser.
//
// It takes three command-line arguments:
//   1) the input obo file
//   2) the nodes output file (3-column tab-delimited)
//   3) the edges output file (3-column tab-delimited, for ancestor relationships)
//
// Note: the parser does not support escape characters except for \" within a
// double-quoted string. Multi-line quoted strings are also not supported.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func parseOBOString(line string) string {
	// parse an optionally quoted string
	line = strings.Replace(line, "\t", " ", -1)
	line = strings.TrimSpace(line)
	if line == "" {
		return line
	}
	if line[0] == '"' {
		// find the next un-escaped " in the line
		idx := 1 + strings.Index(line[1:], `"`)
		for line[idx-1] == '\\' {
			idx = idx + 1 + strings.Index(line[idx+1:], `"`)
		}
		return strings.TrimSpace(line[1:idx])
	}
	return line
}

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "USAGE: %s input.obo nodes.txt edges.txt\n", os.Args[0])
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

	stanza := ""
	id := ""
	name := ""
	def := ""
	edges := make(map[string][]string)

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}
		if line[0] == '[' { // new stanza
			if id != "" && stanza == "[Term]" {
				fmt.Fprintf(nfile, "%s\t%s\t%s\n", id, name, def)
				for rel, others := range edges {
					for _, other := range others {
						fmt.Fprintf(efile, "%s\t%s\t%s\n", id, other, rel)
					}
				}
			}
			stanza = line
			id = ""
			name = ""
			def = ""
			for x := range edges {
				delete(edges, x)
			}
			continue
		}

		row := strings.SplitN(line, ": ", 2)
		if len(row) != 2 {
			panic("tag-value parse error!")
		}
		switch row[0] {
		case "id":
			// spec says id will not contain whitespace
			id = strings.TrimSpace(row[1])

		case "name":
			// spec doesn't say these are quoted, but it doesn't hurt...
			name = parseOBOString(row[1])

		case "def":
			def = parseOBOString(row[1])

		case "is_a":
			parts := strings.Split(row[1], "!")
			// parts[0] = id, so again no whitespace per spec
			edges["is_a"] = append(edges["is_a"], strings.TrimSpace(parts[0]))

		case "relationship":
			parts := strings.Split(row[1], "!")
			parts = strings.Split(parts[0], " ")
			// parts[0] = relationship type
			// parts[1] = id, so again no whitespace per spec
			// parts[2] = modifiers if present (not used)
			edges[parts[0]] = append(edges[parts[0]], strings.TrimSpace(parts[1]))

		case "is_obsolete":
			// per spec, id should always come before is_obsolete,
			// so clearing it here effectively omits it from output
			id = ""

		}
	}

	if id != "" && stanza == "[Term]" {
		fmt.Fprintf(nfile, "%s\t%s\t%s\n", id, name, def)
		for rel, others := range edges {
			for _, other := range others {
				fmt.Fprintf(efile, "%s\t%s\t%s\n", id, other, rel)
			}
		}
	}
}
