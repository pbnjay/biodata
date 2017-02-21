#######################
# transformation tools

tools/parse_obo:
	@mkdir -p tools
	go build -o tools/parse_obo github.com/pbnjay/biodata/parse_obo

tools/parse_mesh:
	@mkdir -p tools
	go build -o tools/parse_mesh github.com/pbnjay/biodata/parse_mesh
