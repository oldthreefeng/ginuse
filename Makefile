.PHONY: release

REGISTRY?=louisehong/deploy
#OLD_TAG?=$(shell git describe --tags | awk -F- '{print $$1}')
OLD_TAG?=$(shell git tag  | sort -n | tail -n 1)
#NEW_TAG=$(shell echo $(OLD_TAG) | awk -F. -v OFS=. 'NF==1{print ++$$NF}; NF>1{if(length($$NF+1)>length($$NF))$$(NF-1)++; $$NF=sprintf("%0*d", length($$NF), ($$NF+1)%(10^length($$NF))); print}')
NEW_TAG=$(shell echo $(OLD_TAG) | awk -F. -v OFS=. 'NF==1{print ++$$NF}; NF>1{if(length($$NF+1)>length($$NF))$$(NF-1)++; $$NF=sprintf("%0*d", length($$NF), ($$NF+1)%(10^length($$NF))); print}')

release:
	-git tag $(NEW_TAG) && git push --tags
	docker build -t  $(REGISTRY):$(NEW_TAG) .
#	docker push $(REGISTRY):$(NEW_TAG)
