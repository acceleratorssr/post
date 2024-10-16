SCRIPTS_DIR := scripts
PYTHON := python3
SHELL := /bin/bash

.PHONY: all clean run_py build_py run_sh build_sh run_ps1 build_ps1

run_py:
	@echo "Running Python script..."
	@cd $(SCRIPTS_DIR) && $(PYTHON) run.py

build_py:
	@echo "Building Python scripts..."
	@cd $(SCRIPTS_DIR) && $(PYTHON) build-all.py

run_sh:
	@echo "Running..."
	@cd $(SCRIPTS_DIR) && ./run.sh

build_sh:
	@echo "Building..."
	@cd $(SCRIPTS_DIR) && ./build-all.sh

run_ps1:
	@echo "Running..."
	@cd $(SCRIPTS_DIR) && powershell -ExecutionPolicy Bypass -File run.ps1

build_ps1:
	@echo "Building..."
	@cd $(SCRIPTS_DIR) && powershell -ExecutionPolicy Bypass -File build-all.ps1

test:

clean:
	@echo "Cleaning up..."
	@cd $(OUTPUT_DIR) && rm -rf *
