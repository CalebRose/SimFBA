package structs

import "testing"

func TestGetCFBStreamFeed(t *testing.T) {
	cases := []struct {
		name       string
		home, away bool
		want       string
	}{
		{"both FBS", true, true, "fbs"},
		{"FBS vs FCS", true, false, "fbs"},
		{"FCS vs FBS", false, true, "fbs"},
		{"both FCS", false, false, "fcs"},
	}
	for _, c := range cases {
		if got := GetCFBStreamFeed(c.home, c.away); got != c.want {
			t.Errorf("%s: got %q want %q", c.name, got, c.want)
		}
	}
}

func TestGetCFBStreamChannel(t *testing.T) {
	cases := []struct {
		timeslot   string
		home, away bool
		want       string
	}{
		{"Thursday Night", true, true, "simfba-live"},
		{"Friday Night", true, false, "simfba-live"},
		{"Saturday Morning", false, true, "simfba-live"},
		{"Saturday Afternoon", true, true, "b1g-noon-kickoff"},
		{"Saturday Evening", true, true, "prime-time"},
		{"Saturday Night", true, true, "pac-2-afterdark"},
		{"", true, true, "simfba-live"},
		{"Saturday Night", false, false, "simfcs-live"},
		{"Thursday Night", false, false, "simfcs-live"},
	}
	for _, c := range cases {
		g := CollegeGame{TimeSlot: c.timeslot}
		if got := GetCFBStreamChannel(g, c.home, c.away); got != c.want {
			t.Errorf("%q home=%v away=%v: got %q want %q", c.timeslot, c.home, c.away, got, c.want)
		}
	}
}

func TestGetNFLStreamChannel(t *testing.T) {
	if got := GetNFLStreamChannel(); got != "simfl-live" {
		t.Errorf("got %q want simfl-live", got)
	}
}
