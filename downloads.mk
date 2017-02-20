
# lists all the files available
DOWNLOADS:=go.obo gene_info.gz gene2go.gz gene2pubmed.gz mesh_d.txt mesh_c.txt mesh_q.txt

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
