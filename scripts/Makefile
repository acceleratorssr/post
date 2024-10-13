SSO_DIR := ../sso
ARTICLE_DIR := ../article
INTERACTIVE_DIR := ../article
USER_DIR := ../user
SEARCH_DIR := ../search
BFF_DIR := ../bff

PYTHON := python3

.PHONY: all clean run_py build_py

sso:
	@echo "Building SSO services.."
	@ $(PYTHON) sso-build.py

art:
	@echo "Building art services.."
	@ $(PYTHON) art-build.py

intr:
	@echo "Building intr services.."
	@ $(PYTHON) intr-build.py

user:
	@echo "Building user services.."
	@ $(PYTHON) user-build.py

search:
	@echo "Building search services.."
	@ $(PYTHON) search-build.py

bff:
	@echo "Building bff services.."
	@ $(PYTHON) bff-build.py