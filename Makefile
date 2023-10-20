MAKEFILE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
ABIS := $(wildcard $(MAKEFILE_DIR)abis/*.abi.json)

.PHONY: list-abis compile

list-abis:
	@echo DIR: $(MAKEFILE_DIR)abis
	@echo ABIs: $(notdir $(ABIS))

compile:
	@echo Generating ABIs...
	@$(foreach abi,$(ABIS), \
		filename=$(notdir $(abi));\
		contractname=$${filename%.*.*}; \
		dir=$(MAKEFILE_DIR)abigen/$${contractname};\
		mkdir -p $${dir};\
		abigen --abi $(abi) --pkg $${contractname} --out $${dir}/$${contractname}.go; \
		echo $${filename};\
	)
