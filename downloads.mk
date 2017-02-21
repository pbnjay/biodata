
# lists all the files available
DOWNLOADS:=go.obo gene_info.gz gene2go.gz gene2pubmed.gz mesh_d.txt mesh_c.txt mesh_q.txt

# used to pick the right MeSH dataset year
YEAR4:=$(shell date +%Y)

fetchall: $(DOWNLOADS)
fetchclean:
	rm -f $(DOWNLOADS)

fetch: fetchclean fetchall

###############
# download files from each of the resources

go.obo:
	curl -LO http://geneontology.org/ontology/go.obo

gene_info.gz:
	curl -LO ftp://ftp.ncbi.nlm.nih.gov/gene/DATA/gene_info.gz

gene2go.gz:
	curl -LO ftp://ftp.ncbi.nlm.nih.gov/gene/DATA/gene2go.gz

gene2pubmed.gz:
	curl -LO ftp://ftp.ncbi.nlm.nih.gov/gene/DATA/gene2pubmed.gz


.mesh_agreement:
	@echo "====="
	@echo
	@echo "If you accept the conditions here: https://www.nlm.nih.gov/mesh/asc_abt.html"
	@echo
	@echo "Then type 'touch .mesh_agreement' and try again."
	@echo
	@exit 1

mesh_d.txt: .mesh_agreement
	curl -L -o mesh_d.txt ftp://nlmpubs.nlm.nih.gov/online/mesh/MESH_FILES/asciimesh/d${YEAR4}.bin

mesh_c.txt: .mesh_agreement
	curl -L -o mesh_c.txt ftp://nlmpubs.nlm.nih.gov/online/mesh/MESH_FILES/asciimesh/c${YEAR4}.bin

mesh_q.txt: .mesh_agreement
	curl -L -o mesh_q.txt ftp://nlmpubs.nlm.nih.gov/online/mesh/MESH_FILES/asciimesh/q${YEAR4}.bin
