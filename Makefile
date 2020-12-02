run-all:
	for d in day-*; \
		do echo "=== $$(echo $$d | tr a-z A-Z) ==="; \
		go run ./$$d < $$d/input.txt; \
		echo; \
	done

.PHONY: run-all