YEAR ?= $(shell date "+%Y")

run-all:
	for d in puzzles/$(YEAR)/day-*; \
		do echo "=== $$(echo $$(basename $$d) | tr a-z A-Z) ==="; \
		go run ./$$d < $$d/input.txt; \
		echo; \
	done

leaderboard:
	go run admin.go leaderboard --id $(AOC_PRIVATE_LEADERBOARD_ID) 

init-next:
	go run admin.go bootstrap

.PHONY: run-all leaderboard