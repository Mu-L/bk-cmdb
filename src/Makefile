
IMAGE?=image_placeholder
NPM?=npm
PROJECT_PATH=$(shell cd ../;  pwd)

include ../scripts/Makefile

.PHONY:default
default:
	@echo ""
	@cd  $(SCRIPT_DIR) && DISABLE_CRYPTO=$(DISABLE_CRYPTO) bash ./tongsuo.sh
	@echo -e "\033[34mbuilding server... \033[0m"
	@rm -rf "./ui/node_modules"
	@cd  $(SCRIPT_DIR) && bash ./build.sh
	@cd $(SCRIPT_DIR) && python ./generate.py -t '$(BIN_PATH)' -i '${IMAGE}' -o '$(BIN_PATH)/docker'
	@mkdir $(BIN_PATH)/changelog_user/
	@cp -R ${DOCS_DIR}/support-file/changelog_user/*  $(BIN_PATH)/changelog_user
	@echo -e "\033[34mbuilding web... \033[0m"
	@cd $(SOURCE_ROOT)/ui && bash ./build.sh ${NPM} ${BIN_PATH}

.PHONY:linux
linux:
	@echo ""
	@cd  $(SCRIPT_DIR) && DISABLE_CRYPTO=$(DISABLE_CRYPTO) bash ./tongsuo.sh
	@echo -e "\033[34mbuilding server... \033[0m"
	@rm -rf "./ui/node_modules"
	@cd  $(SCRIPT_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 bash ./build.sh
	@cd $(SCRIPT_DIR) && python ./generate.py -t '$(BIN_PATH)' -i '${IMAGE}' -o '$(BIN_PATH)/docker'
	@mkdir $(BIN_PATH)/changelog_user/
	@cp -R ${DOCS_DIR}/support-file/changelog_user/*  $(BIN_PATH)/changelog_user
	@echo -e "\033[34mbuilding web... \033[0m"
	@cd $(SOURCE_ROOT)/ui && bash ./build.sh ${NPM} ${BIN_PATH}

.PHONY:server
server:
	@echo ""
	@cd  $(SCRIPT_DIR) && DISABLE_CRYPTO=$(DISABLE_CRYPTO) bash ./tongsuo.sh
	@echo -e "\033[34mbuilding server... \033[0m"
	@rm -rf "./ui/node_modules"
	@cd  $(SCRIPT_DIR) &&  bash ./build.sh
	@cd $(SCRIPT_DIR) && python ./generate.py -t '$(BIN_PATH)' -i '${IMAGE}' -o '$(BIN_PATH)/docker'
	@mkdir $(BIN_PATH)/changelog_user/
	@cp -R ${DOCS_DIR}/support-file/changelog_user/*  $(BIN_PATH)/changelog_user

.PHONY:ui
ui:
	@echo ""
	@echo -e "\033[34mbuilding web... \033[0m"
	@rm -rf "./ui/node_modules"
	@cd $(SOURCE_ROOT)/ui && bash ./build.sh ${NPM} ${BIN_PATH}

.PHONY:debug
debug:
	@echo ""
	@echo -e "\033[34mbuilding debug... \033[0m"
	@rm -rf "./ui/node_modules"
	@cd  $(SCRIPT_DIR) && DISABLE_CRYPTO=$(DISABLE_CRYPTO) bash ./tongsuo.sh
	@cd $(SCRIPT_DIR) && bash ./build.sh debug
	@cd $(SCRIPT_DIR) && python ./generate.py -t '$(BIN_PATH)' -i '${IMAGE}' -o '$(BIN_PATH)/docker'

.PHONY:package
package:
	rm -rf $(SOURCE_ROOT)/bin/pub/cmdb
	mkdir -p $(SOURCE_ROOT)/bin/pub && cp -R $(BIN_PATH) ./bin/pub/cmdb
	cd $(SOURCE_ROOT)/bin/pub && tar -zvcf cmdb.tar.gz cmdb
	@echo -e "\033[34mpackage success and store in $(SOURCE_ROOT)/bin/pub/cmdb.tar.gz\033[0m"

.PHONY:enterprise
enterprise:
	rm -rf $(SOURCE_ROOT)/bin/enterprise/cmdb
	mkdir -p $(SOURCE_ROOT)/bin/enterprise/cmdb/server/bin/
	mkdir -p $(SOURCE_ROOT)/bin/enterprise/cmdb/errors/
	mkdir -p $(SOURCE_ROOT)/bin/enterprise/cmdb/language/
	mkdir -p $(SOURCE_ROOT)/bin/enterprise/cmdb/web/
	mkdir -p $(SOURCE_ROOT)/bin/enterprise/cmdb/monstache/etc
	mkdir -p $(SOURCE_ROOT)/bin/enterprise/cmdb/support-files/
	mkdir -p $(SOURCE_ROOT)/bin/enterprise/cmdb/server/changelog_user/
	mkdir -p $(SOURCE_ROOT)/bin/enterprise/cmdb/apidocs/apidocs/

	cp -R ${RESOURCE_DIR}/errors/*  $(SOURCE_ROOT)/bin/enterprise/cmdb/errors
	cp -R ${RESOURCE_DIR}/language/*  $(SOURCE_ROOT)/bin/enterprise/cmdb/language
	cp -R ${BIN_PATH}/monstache/  $(SOURCE_ROOT)/bin/enterprise/cmdb/monstache
	cp -R ${BIN_PATH}/web/*  $(SOURCE_ROOT)/bin/enterprise/cmdb/web/
	cp  ${BIN_PATH}/cmdb_*/cmdb_*  $(SOURCE_ROOT)/bin/enterprise/cmdb/server/bin/
	cp  ${BIN_PATH}/tool_*/tool_*  $(SOURCE_ROOT)/bin/enterprise/cmdb/server/bin/
	mv ${DOCS_DIR}/support-file/config/projects.yaml $(SOURCE_ROOT)/bin/enterprise/cmdb/projects.yaml
	cp -R ${DOCS_DIR}/support-file/config/*  $(SOURCE_ROOT)/bin/enterprise/cmdb/support-files
	cp -R ${DOCS_DIR}/support-file/changelog/*  $(SOURCE_ROOT)/bin/enterprise/cmdb
	cp -R ${DOCS_DIR}/support-file/changelog_user/*  $(SOURCE_ROOT)/bin/enterprise/cmdb/server/changelog_user
	cp -R ${DOCS_DIR}/apidoc/apigw/definition.yaml  $(SOURCE_ROOT)/bin/enterprise/cmdb/apidocs/definition.yaml
	cp -R ${DOCS_DIR}/apidoc/apigw/backend/bk_apigw_resources_bk-cmdb.yaml  $(SOURCE_ROOT)/bin/enterprise/cmdb/apidocs/resources.yaml
	cp -R ${DOCS_DIR}/apidoc/apigw/bin  $(SOURCE_ROOT)/bin/enterprise/cmdb/apidocs
	cp -R ${DOCS_DIR}/apidoc/apigw/backend/en  $(SOURCE_ROOT)/bin/enterprise/cmdb/apidocs/apidocs
	cp -R ${DOCS_DIR}/apidoc/apigw/backend/zh  $(SOURCE_ROOT)/bin/enterprise/cmdb/apidocs/apidocs
	cp -R ${DOCS_DIR}/apidoc/apigw/open/zh/* $(SOURCE_ROOT)/bin/enterprise/cmdb/apidocs/apidocs/zh
	cp -R ${DOCS_DIR}/apidoc/apigw/open/en/* $(SOURCE_ROOT)/bin/enterprise/cmdb/apidocs/apidocs/en
	tail -n +10 ${DOCS_DIR}/apidoc/apigw/open/bk_apigw_resources_bk-cmdb.yaml >> $(SOURCE_ROOT)/bin/enterprise/cmdb/apidocs/resources.yaml
	for file in ${DOCS_DIR}/apidoc/apigw/open/zh/*; do \
		cp $$file $(SOURCE_ROOT)/bin/enterprise/cmdb/apidocs/apidocs/zh/open_$$(basename $$file); \
	done
	for file in ${DOCS_DIR}/apidoc/apigw/open/en/*; do \
		cp $$file $(SOURCE_ROOT)/bin/enterprise/cmdb/apidocs/apidocs/en/open_$$(basename $$file); \
	done
	tail -n +10 ${DOCS_DIR}/apidoc/apigw/open/bk_apigw_resources_bk-cmdb.yaml | \
	sed -e 's|  /api/v3/|  /api/v3/open/|g' \
		-e 's|isPublic: false|isPublic: true|g' \
		-e 's|userVerifiedRequired: false|userVerifiedRequired: true|g' \
		-e 's|operationId: |operationId: open_|g' >> $(SOURCE_ROOT)/bin/enterprise/cmdb/apidocs/resources.yaml

.PHONY:clean
clean:
	rm -rf ${BIN_PATH}

.PHONY:cleanall
cleanall:
	rm -rf $(SOURCE_ROOT)/bin
