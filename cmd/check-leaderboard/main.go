// check-leaderboard will get the lastest state of a private leaderboard,
// given a leaderboard ID and a session token.
//
// If no ID and/or no session token are provided, it will decode a JSON
// blob that represents a leaderboard from stdin.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Determine format for leaderboard summary: how to show progress by day.
// Track last-fetch time and last response, to keep from over-checking.

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		leaderboardID uint
		sessionCookie string
	)
	flag.UintVar(&leaderboardID, "id", 0, "private leaderboard ID")
	flag.StringVar(&sessionCookie, "session", "", "session cookie value")
	flag.Parse()

	var (
		reader = ioutil.NopCloser(os.Stdin)
		err    error
	)

	// Fetch from internet if ID and session spcified.
	if leaderboardID > 0 && sessionCookie != "" {
		reader, err = fetchLeaderboard(http.DefaultClient, leaderboardID, sessionCookie)
		if err != nil {
			return fmt.Errorf("fetching leaderboard: %w", err)
		}
	}

	var lb Leaderboard
	if err := json.NewDecoder(reader).Decode(&lb); err != nil {
		return fmt.Errorf("decoding leaderboard: %w", err)
	}

	fmt.Println(lb.Summary())
	return nil
}

func fetchLeaderboard(client *http.Client, id uint, sessionCookie string) (io.ReadCloser, error) {
	lbURL := fmt.Sprintf("https://adventofcode.com/2020/leaderboard/private/view/%d.json", id)
	r, err := http.NewRequest(http.MethodGet, lbURL, nil)
	if err != nil {
		return nil, err
	}
	r.AddCookie(&http.Cookie{Name: "session", Value: sessionCookie})

	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

type (
	// Leaderboard contains the top level data for the leaderboard and
	// all the stats for members.
	Leaderboard struct {
		OwnerID string            `json:"owner_id"`
		Event   string            `json:"event"`
		Members map[string]Member `json:"members"`
	}

	// Member describes an individual member of the leaderboard, including
	// identity and leaderboard stats.
	Member struct {
		ID                 string             `json:"id"`
		Name               string             `json:"name"`
		Stars              int                `json:"stars"`
		GlobalScore        int                `json:"global_score"`
		LocalScore         int                `json:"local_score"`
		LastStarTimestamp  time.Time          `json:"last_star_ts,omitempty"`
		CompletionDayLevel map[int]DailyStats `json:"completion_day_level"`
	}

	// DailyStats contains all the stats for a given day of an event.
	DailyStats map[int]StarTimestamp

	// StarTimestamp indicates when a star was awarded.
	StarTimestamp struct {
		GetStarTimestamp time.Time `json:"get_star_ts"`
	}
)

func (lb Leaderboard) Summary() string {
	const (
		headingPoints         = "POINTS"
		headingMostRecentStar = "MOST RECENT STAR"
		dateFormat            = "2006-01-02 15:04:05 -0700 MST"
	)
	var (
		longest      = lb.LongestMemberNameLen()
		headerFormat = fmt.Sprintf("%%-%ds  %s  %s\n", longest, headingPoints, headingMostRecentStar)
		underline    = strings.Repeat("=", longest) + "  " +
			strings.Repeat("=", len(headingPoints)) + "  " +
			strings.Repeat("=", len(dateFormat)) + "\n"
		entryFormat = fmt.Sprintf(
			"%%-%ds  %%%dd  %%%ds\n",
			longest,
			len(headingPoints),
			len(dateFormat))
	)
	summary := fmt.Sprintf(headerFormat, "NAME")
	summary += underline
	for _, m := range lb.SortedMembers() {
		lastStarTimestamp := fmt.Sprint(m.LastStarTimestamp)
		if m.LastStarTimestamp.Unix() == 0 || m.LastStarTimestamp.IsZero() {
			lastStarTimestamp = "(none)"
		}
		summary += fmt.Sprintf(entryFormat, m.Name, m.LocalScore, lastStarTimestamp)
	}
	return summary
}

// LongestNameLen gets the length of the longest member name.
func (lb Leaderboard) LongestMemberNameLen() int {
	max := 0
	for _, v := range lb.Members {
		if l := len(v.Name); l > max {
			max = l
		}
	}
	return max
}

func (lb Leaderboard) SortedMembers() []Member {
	members := make([]Member, 0, len(lb.Members))
	for _, m := range lb.Members {
		members = append(members, m)
	}
	sort.Slice(members, func(i, j int) bool {
		// Reverse the definition of a typical "less" here so that the sort
		// comes back in reverse order, with highest value first.
		return members[i].LocalScore > members[j].LocalScore
	})
	return members
}

func (m *Member) UnmarshalJSON(b []byte) error {
	var member struct {
		ID                 string             `json:"id"`
		Name               string             `json:"name"`
		Stars              int                `json:"stars"`
		GlobalScore        int                `json:"global_score"`
		LocalScore         int                `json:"local_score"`
		LastStarTimestamp  json.Number        `json:"last_star_ts"`
		CompletionDayLevel map[int]DailyStats `json:"completion_day_level"`
	}

	if err := json.Unmarshal(b, &member); err != nil {
		return err
	}

	m.ID = member.ID
	m.Name = member.Name
	m.Stars = member.Stars
	m.GlobalScore = member.GlobalScore
	m.LocalScore = member.LocalScore
	m.CompletionDayLevel = member.CompletionDayLevel

	v, err := member.LastStarTimestamp.Int64()
	if err != nil {
		return err
	}
	m.LastStarTimestamp = time.Unix(v, 0)

	return nil
}

func (st *StarTimestamp) UnmarshalJSON(b []byte) error {
	var s struct {
		Timestamp json.Number `json:"get_star_ts"`
	}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	ts, err := s.Timestamp.Int64()
	if err != nil {
		return err
	}
	st.GetStarTimestamp = time.Unix(ts, 0)
	return nil
}

func (st StarTimestamp) MarshalJSON() ([]byte, error) {
	ts := st.GetStarTimestamp.Unix()
	out := struct {
		GetStarTimestamp string `json:"get_star_ts"`
	}{
		GetStarTimestamp: strconv.FormatInt(ts, 10),
	}
	return json.Marshal(out)
}
