#######################
# transformation tools

tools/parse_obo:
	@mkdir -p tools
	go build -o tools/parse_obo github.com/pbnjay/biodata/parse_obo
