# Переменные
GOOS=linux
GOARCH=amd64
BUILD_DIR=build
FUNCTIONS=image-analyzer translator speech product-search

# AWS переменные
AWS_REGION=eu-north-1
ACCOUNT_ID=730918949094

# Цвета для вывода
GREEN=\033[0;32m
NC=\033[0m # No Color

.PHONY: all clean build package $(FUNCTIONS) deploy-% build-% package-%

# Создаем все Lambda функции
all: clean build package

# Очистка build директории
clean:
	@echo "${GREEN}Cleaning build directory...${NC}"
	@rm -rf $(BUILD_DIR)
	@mkdir -p $(BUILD_DIR)
	@for func in $(FUNCTIONS); do \
		mkdir -p $(BUILD_DIR)/$$func; \
	done

# Сборка всех функций
build: $(FUNCTIONS)

# Паттерн для сборки каждой функции
$(FUNCTIONS):
	@echo "${GREEN}Building $@ lambda...${NC}"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BUILD_DIR)/$@/bootstrap cmd/$@/main.go

# Упаковка в ZIP
package: $(FUNCTIONS:%=%-zip)

# Паттерн для создания ZIP архива
%-zip:
	@echo "${GREEN}Creating ZIP for $*...${NC}"
	@cd $(BUILD_DIR)/$* && zip -q ../$*.zip bootstrap
	@echo "${GREEN}Created $(BUILD_DIR)/$*.zip${NC}"

# Команды для отдельных функций
build-%:
	@echo "${GREEN}Building $* lambda...${NC}"
	@mkdir -p $(BUILD_DIR)/$*
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BUILD_DIR)/$*/bootstrap cmd/$*/main.go

package-%:
	@echo "${GREEN}Creating ZIP for $*...${NC}"
	@cd $(BUILD_DIR)/$* && zip -q ../$*.zip bootstrap
	@echo "${GREEN}Created $(BUILD_DIR)/$*.zip${NC}"

deploy-%:
	@echo "${GREEN}Deploying $* lambda...${NC}"
	@aws lambda update-function-code \
		--region $(AWS_REGION) \
		--function-name $* \
		--zip-file fileb://$(BUILD_DIR)/$*.zip \
		--publish
	@echo "${GREEN}Deployed $* lambda${NC}"

# Команды для каждой функции
.PHONY: translator-all
translator-all: build-translator package-translator deploy-translator

.PHONY: image-analyzer-all
image-analyzer-all: build-image-analyzer package-image-analyzer deploy-image-analyzer

.PHONY: speech-all
speech-all: build-speech package-speech deploy-speech

.PHONY: product-search-all
product-search-all: build-product-search package-product-search deploy-product-search

# Показать список доступных команд
help:
	@echo "Available commands:"
	@echo "  make all                    - Clean, build and package all functions"
	@echo "  make clean                  - Clean build directory"
	@echo "  make build                  - Build all functions"
	@echo "  make package                - Package all functions into ZIP files"
	@echo "  make build-<function>       - Build specific function (e.g., make build-translator)"
	@echo "  make package-<function>     - Package specific function (e.g., make package-translator)"
	@echo "  make deploy-<function>      - Deploy specific function (e.g., make deploy-translator)"
	@echo "  make <function>-all         - Build, package and deploy specific function (e.g., make translator-all)"
	@echo ""
	@echo "Available functions:"
	@echo "  - translator"
	@echo "  - image-analyzer"
	@echo "  - speech"
	@echo "  - product-search"
	@echo ""
	@echo "Examples:"
	@echo "  make translator-all         - Build, package and deploy translator function"
	@echo "  make build-speech          - Only build speech function"
	@echo "  make deploy-image-analyzer - Only deploy image-analyzer function"

# Тестирование сборки одной функции
test-build:
	@echo "${GREEN}Testing build process...${NC}"
	@make clean
	@make image-analyzer
	@make image-analyzer-zip
	@echo "${GREEN}Test build completed. Check build/image-analyzer.zip${NC}" 