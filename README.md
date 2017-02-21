# biodata
Scripts and tools for downloading and parsing public datasets for bioinformatics.
This is particularly oriented to knowledge graphs and data integration projects.

# Intentions
If you are new to ETL / bioinformatics, then this repository will give you many
examples and working code and allow you to get up and running quickly with
knowledge-oriented queries. If you have experience already, this repo hopefully
provides a good baseline for datasets that are new to you.

This repository is intended as a "starting off" point for more complex usage.
If you need more data than what is currently provided, you should learn more
about the data source itself and modify any code found here to do what you need!

# Quick Start
Make sure you have [Go installed](https://golang.org/dl/) along with GNU Make,
then clone this repo, cd to it and run `make`. This will download and process
all the data files and produce partitioned data as described below. If you would
like to play with the code and other scripts you should clone this into the
appropriate GOPATH.

# Data Model
Data processing currently uses a "partitioned graph" conceptual model focused
on "nodes" within "partitions", and "edges" between nodes. This is reasonably
straightforward to apply to a relational model of tables and relations.

  - Partition: provides a grouping or namespace for sets of nodes:
    - In this repository, we generally use a reverse-domain mapping for the
      resource's descriptive site URL as the partition name. In some cases, a
      defined identifier may already exist (i.e. found within RDF/OWL schemas).
    - A public resource will often have multiple partitions of data:
      - `gov.nih.nlm.ncbi.gene`: [NCBI Gene database](https://www.ncbi.nlm.nih.gov/gene)
      - `gov.nih.nlm.ncbi.taxonomy`: [NCBI Taxonomy database](https://www.ncbi.nlm.nih.gov/taxonomy)
      - `gov.nih.nlm.ncbi.pubmed`: [NCBI Pubmed database](https://www.ncbi.nlm.nih.gov/pubmed)
  
  - Nodes: three pieces of basic information about the data point:
    - Identifier - a value that uniquely identifies the data within the partition.
    - Name - a text value that is more human-friendly for display.
    - Description - a longer text value with more information.

  - Edges: relate two nodes to each other:
    - "Left" and "Right" nodes, with a relation type connecting them.
    - Relation types are specified as from left to right. e.g. "left is part of right"
    - Relations may have a value associated. e.g. "left is correlated at 0.71 to right"
    - If you're familiar with the terms, these are basically [Triples](https://en.wikipedia.org/wiki/Semantic_triple) (Quads) per RDF

# Outputs
Output of these tools are simple tab-delimited text files, organized into
directories named by each partition. Results are typically gzipped to conserve
storage space (all of the bundled tools will transparently handle gzip files
as input or output).

  - Node definitions are always found in a file named `nodes.txt` (e.g. found in `dots.to.partition/nodes.txt`).
  	- This file is 3 tab-delimited columns `identifier <TAB> name <TAB> description` with no header.
  - Edges are found in the "left" partition's directory, named by the "right"
  	paritition, e.g. `dots.to.left.partition/dots.to.right.partition.txt`.
  	- These files have 3-4 columns `left-id <TAB> right-id <TAB> relation [<TAB> relation-value]` with no header.
  	- Note the column ordering is intentional, as it allows code to easily ignore
  	  the semantic types of relations if desired, and groups possible values to
  	  the relation more closely.

# Tools
This repository will try to stick to standard unix command line tools as much as
possible. If necessary, stand-alone tools will be added to aid in the extraction
of the types of data listed above. All tools will be written in [the Go programming language](https://golang.org/)
to ensure portability and speed. Once you have a Go toolchain installed, these
tools can be easily installed with: `go get github.com/pbnjay/biodata/<tool_name>`

# Contributing
All of the resources here should be publicly discoverable from a stable provider.
Commercial resources will not be included within this project (non-profits that
require payment for licensing/API access, etc. are allowed, but will require
user-provided authentication to operate).