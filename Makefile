# miscellaneous definition used by dependencies
include defs.mk

ALLPARTS=go gene mesh
all: $(ALLPARTS)

# defines how to build tools included with this repo
include tools.mk

# defines how to download the files needed for processing
include downloads.mk 

###############
# dataset production

go: tools/parse_obo go.obo
	@mkdir -p ${PART_GO}
	tools/parse_obo go.obo ${PART_GO}/nodes.txt ${PART_GO}/${PART_GO}.txt
	sort -u ${PART_GO}/${PART_GO}.txt | gzip -9 > ${PART_GO}/${PART_GO}.txt.gz
	sort -u ${PART_GO}/nodes.txt | gzip -9 > ${PART_GO}/nodes.txt.gz
	rm ${PART_GO}/${PART_GO}.txt ${PART_GO}/nodes.txt

gene: gene_go gene_nodes gene_pubmed

gene_go: gene2go.gz
	@mkdir -p ${PART_GENE}
	gzcat gene2go.gz | awk -F "\t" '{if(NR>1) {print $$2 "\t" $$3 "\tgo:" $$4}}' | sort -u > ${PART_GENE}/${PART_GO}.txt
	gzip -9 ${PART_GENE}/${PART_GO}.txt

gene_nodes: gene_info.gz
	@mkdir -p ${PART_GENE}
	gzcat gene_info.gz | awk -F "\t" '{if(NR>1) {print $$2 "\t" $$3 "\t" $$9}}' | sort -u > ${PART_GENE}/nodes.txt
	gzip -9 ${PART_GENE}/nodes.txt

gene_pubmed: gene2pubmed.gz
	@mkdir -p ${PART_GENE}
	gzcat gene2pubmed.gz | awk -F "\t" '{if(NR>1) {print $$2 "\t" $$3 "\treferenced_in"}}' | sort -u > ${PART_GENE}/${PART_PUBMED}.txt
	gzip -9 ${PART_GENE}/${PART_PUBMED}.txt

mesh: tools/parse_mesh mesh_c.txt mesh_d.txt mesh_q.txt
	@mkdir -p ${PART_MESH}
	tools/parse_mesh mesh_d.txt ${PART_MESH}/nodes.d.txt ${PART_MESH}/${PART_MESH}.d.txt
	# these don't actually have edges, but we try anyway...
	tools/parse_mesh mesh_c.txt ${PART_MESH}/nodes.c.txt ${PART_MESH}/${PART_MESH}.c.txt
	tools/parse_mesh mesh_q.txt ${PART_MESH}/nodes.q.txt ${PART_MESH}/${PART_MESH}.q.txt
	# combine all three!
	sort -u ${PART_MESH}/${PART_MESH}.?.txt | gzip -9 > ${PART_MESH}/${PART_MESH}.txt.gz
	sort -u ${PART_MESH}/nodes.?.txt | gzip -9 > ${PART_MESH}/nodes.txt.gz
	rm ${PART_MESH}/${PART_MESH}.?.txt ${PART_MESH}/nodes.?.txt

pubmed: tools/parse_pubmed medline_files.txt
	@mkdir -p ${PART_PUBMED}
	tools/parse_pubmed pubmed_data ${PART_PUBMED}/nodes.txt.gz ${PART_PUBMED}/${PART_MESH}.txt.gz

################

# Three ways to skip 1 header line:
#  'tail -n +2'
#   - can be slow using BSD/macOS tail, fine on linux
#
#  "awk '{if(NR>1)..."
#   - good esp when you need to mix prefixes or multi-column extractions
#
#  'sed 1d'
#   - BSD/macOS sed is different enough from linux to make it annoying beyond simple usage
