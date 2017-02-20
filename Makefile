# miscellaneous definition used by dependencies
include defs.mk

# defines how to build tools included with this repo
include tools.mk

# defines how to download the files needed for processing
include downloads.mk 

###############
# dataset production

# this makes sorts much faster by using native byte order...
export LC_ALL=C

# Three ways to skip 1 header line:
#  'tail -n +2'
#   - can be slow using BSD/macOS tail, fine on linux
#
#  "awk '{if(NR>1)..."
#   - good esp when you need to mix prefixes or multi-column extractions
#
#  'sed 1d'
#   - BSD/macOS sed is different enough from linux to make it annoying beyond simple usage

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
