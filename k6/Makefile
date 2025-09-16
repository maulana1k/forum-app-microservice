# Makefile for K6 load testing with web dashboard

.PHONY: test report clean open

# Configuration
TEST_SCRIPT = stress_test.js
REPORT_FILE = k6-report.html
PORT = 5665

# Run test with web dashboard
test:
	K6_WEB_DASHBOARD=true K6_WEB_DASHBOARD_EXPORT=$(REPORT_FILE) k6 run $(TEST_SCRIPT)

dashboard:
	open http://localhost:$(PORT) && open $(REPORT_FILE)	

# Run test and auto-open report
report: test
	@echo "📊 Opening dashboard..."
	@if command -v xdg-open > /dev/null; then \
		xdg-open http://localhost:$(PORT) && xdg-open $(REPORT_FILE); \
	elif command -v open > /dev/null; then \
		open http://localhost:$(PORT) && open $(REPORT_FILE); \
	else \
		echo "Open manually: http://localhost:$(PORT) and $(REPORT_FILE)"; \
	fi

# Run test in background and open dashboard
live:
	@echo "Starting live dashboard..."
	K6_WEB_DASHBOARD=true k6 run $(TEST_SCRIPT) & \
	PID=$$!; \
	sleep 3; \
	if command -v xdg-open > /dev/null; then \
		xdg-open http://localhost:$(PORT); \
	elif command -v start > /dev/null; then \
		start http://localhost:$(PORT); \
	fi; \
	wait $$PID

# Generate HTML report only
html:
	K6_WEB_DASHBOARD_EXPORT=$(REPORT_FILE) k6 run $(TEST_SCRIPT)
	@echo "📁 Report saved: $(REPORT_FILE)"

# Open existing report
open:
	@if command -v xdg-open > /dev/null; then \
		xdg-open $(REPORT_FILE); \
	elif command -v open > /dev/null; then \
		open $(REPORT_FILE); \
	else \
		echo "Open manually: $(REPORT_FILE)"; \
	fi

# Clean generated files
clean:
	rm -f $(REPORT_FILE) results.json

# Help message
help:
	@echo "Available commands:"
	@echo "  make test      - Run test with web dashboard"
	@echo "  make report    - Run test and open HTML report"
	@echo "  make live      - Run test with live web dashboard"
	@echo "  make html      - Generate HTML report only"
	@echo "  make open      - Open existing report"
	@echo "  make clean     - Clean generated files"