run-all:
	for d in puzzles/day-*; \
		do echo "=== $$(echo $$(basename $$d) | tr a-z A-Z) ==="; \
		go run ./$$d < $$d/input.txt; \
		echo; \
	done

leaderboard:
	go run ./cmd/check-leaderboard \
		--id $(AOC_PRIVATE_LEADERBOARD_ID) \
		--session $(AOC_SESSION_COOKIE)

.PHONY: run-all leaderboard