.PHONY: new-project help

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  new-project  Create a new Go project from a template"
	@echo "  clean        Remove created projects (if any)"

new-project:
	@chmod +x setup.sh
	@./setup.sh
