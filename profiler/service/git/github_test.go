package git

import (
	"encoding/json"
	"testing"
)

func TestFetchUserDetails(t *testing.T) {
	t.Run("Check user details", func(t *testing.T) {
		gs := NewGitService()
		user, err := gs.FetchUserDetails("xhermitx")
		s, _ := json.MarshalIndent(user, "  ", "\t")
		t.Log(s) // Only runs in verbose mode
		if err != nil {
			t.Errorf("error occured: %v", err)
		}
		if user.Username != "xhermitx" {
			t.Error("wrong username")
		}
	})
}
