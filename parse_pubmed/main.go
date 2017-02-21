// Command parse_pubmed will parse pubmed/medline XML collections.
//
// It takes three command-line arguments:
//   1) the root directory containing medline*.xml.gz files
//   2) the nodes output file (3-column tab-delimited)
//   3) the edges output file (4-column tab-delimited, for mesh descriptor/qualifier tags)
//
// Note: many fields within the XML documents are not parsed or used. Updates are not yet
// consolidated so changes to prior data will not be visible.
package main

import (
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type CitationSet struct {
	Citations []Citation `xml:"PubmedArticle>MedlineCitation"`
}

type MeshDescriptor struct {
	MajorTopic string `xml:"MajorTopicYN,attr"`
	MeshUI     string `xml:"UI,attr"`
	Name       string `xml:",chardata"`
}

type MeshQualifiedDescriptor struct {
	Desc MeshDescriptor   `xml:"DescriptorName"`
	Qual []MeshDescriptor `xml:"QualifierName"`
}

type ELocID struct {
	IDType   string `xml:"EIdType,attr"`
	Location string `xml:",chardata"`
}

type Citation struct {
	PMID        int      `xml:"PMID"`
	PubYear     int      `xml:"Article>Journal>JournalIssue>PubDate>Year"`
	PubMonth    string   `xml:"Article>Journal>JournalIssue>PubDate>Month"`
	AltPubYear  int      `xml:"DateCreated>Year"`
	AltPubMonth int      `xml:"DateCreated>Month"`
	Journal     string   `xml:"Article>Journal>Title"`
	JournalAbbr string   `xml:"Article>Journal>ISOAbbreviation"`
	Title       string   `xml:"Article>ArticleTitle"`
	Authors     []string `xml:"Article>AuthorList>Author>LastName"`
	ELocation   ELocID   `xml:"Article>ELocationID"`

	MeshHeadings []MeshQualifiedDescriptor `xml:"MeshHeadingList>MeshHeading"`
}

func parsePubmed(r io.Reader, nfile, efile io.Writer) error {
	var cs CitationSet
	err := xml.NewDecoder(r).Decode(&cs)
	if err != nil {
		return err
	}

	for _, c := range cs.Citations {
		name := ""
		if len(c.Authors) > 0 {
			name = c.Authors[0]
			if len(c.Authors) > 1 {
				name += " et al"
			}
			name += ". "
		}
		pubmon := c.PubMonth
		if pubmon == "" {
			// alternate month is ill-defined, but I think this'll work
			pubmon = time.Month(c.AltPubMonth).String()
		}
		if c.PubYear <= 0 {
			c.PubYear = c.AltPubYear
		}
		pubyear := ""
		if c.PubYear < 100 {
			pubyear = fmt.Sprintf("20%d", c.PubYear)
		} else {
			pubyear = fmt.Sprintf("%d", c.PubYear)
		}

		// if month is empty try to consolidate space
		pubdt := strings.TrimSpace(pubmon + " " + pubyear)

		// don't allow tabs/newlines within titles
		c.Title = strings.Replace(c.Title, "\t", " ", -1)
		c.Title = strings.Replace(c.Title, "\n", " ", -1)

		// sometimes title/journal have dots at the end. standardize them...
		c.Title = strings.TrimRight(c.Title, ". ")
		c.JournalAbbr = strings.TrimRight(c.JournalAbbr, ". ")

		// all the work comes together to build a simple "citation" for the publication description
		// (we don't want to keep abstracts in our nodes...)
		fullcitation := fmt.Sprintf("%s%s. %s. (%s)", name, c.Title, c.JournalAbbr, pubdt)
		fmt.Fprintf(nfile, "%d\t%s\t%s\n", c.PMID, c.Title, fullcitation)

		for _, qd := range c.MeshHeadings {
			h := qd.Desc
			if h.MeshUI == "" {
				panic("no UI for mesh topic?!? " + h.Name)
			}
			ctx := "minor"
			if h.MajorTopic == "Y" {
				ctx = "major"
			}
			// PMID has topic descriptor
			fmt.Fprintf(efile, "%d\t%s\t%s\t%s\n", c.PMID, h.MeshUI, "has_topic", ctx)

			for _, q := range qd.Qual {
				pred := "has_topic_qualifier"
				if h.MajorTopic == "Y" {
					pred = "has_major_topic_qualifier"
				}
				if q.MeshUI == "" {
					panic("no UI for mesh topic?!? " + q.Name)
				}

				// NB this may appear flipped at first glance. The qualifier is primary
				// target, with descriptor specifying context. Note that above edge has
				// descriptor as target so the PMID would be found, but this arrangement
				// allows for more general queries. e.g. Disease X will hit above, but
				// "complications" will hit here, and Disease X complications is
				// possible by querying on both.
				fmt.Fprintf(efile, "%d\t%s\t%s\t%s\n", c.PMID, q.MeshUI, pred, h.MeshUI)
			}
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "USAGE: %s pubmed_dir nodes.txt[.gz] edges.txt[.gz]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  where pubmed_dir is a directory containing medline*.xml.gz files\n")
		fmt.Fprintf(os.Stderr, "    and nodes.txt[.gz] will receive the node listing\n")
		fmt.Fprintf(os.Stderr, "    and edges.txt[.gz] will receive the pubmed->mesh edges listing\n")
		return
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

	var zn io.Writer = nfile
	var ze io.Writer = efile
	if strings.HasSuffix(os.Args[2], ".gz") {
		znf, err := gzip.NewWriterLevel(nfile, gzip.BestCompression)
		if err != nil {
			panic(err)
		}
		defer znf.Close()
		zn = znf
	}

	if strings.HasSuffix(os.Args[3], ".gz") {
		zef, err := gzip.NewWriterLevel(efile, gzip.BestCompression)
		if err != nil {
			panic(err)
		}
		defer zef.Close()
		ze = zef
	}

	err = filepath.Walk(os.Args[1], func(wpath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !strings.HasPrefix(info.Name(), "medline") || !strings.HasSuffix(info.Name(), ".xml.gz") {
			return nil
		}

		f, err := os.Open(wpath)
		if err != nil {
			return err
		}
		defer f.Close()

		rz, err := gzip.NewReader(f)
		if err != nil {
			return err
		}
		defer rz.Close()

		log.Println(wpath)
		////////////////
		return parsePubmed(rz, zn, ze)
	})

	if err != nil {
		panic(err)
	}
}
