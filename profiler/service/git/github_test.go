package git

import (
	"testing"
)

func TestFetchUserDetails(t *testing.T) {
	t.Run("Check user details", func(t *testing.T) {
		gs := NewGitService()
		user, err := gs.FetchUserDetails("xhermitx")
		t.Log(user) // Only runs in verbose mode
		if err != nil {
			t.Errorf("error occured: %v", err)
		}
		if user.Username != "xhermitx" {
			t.Error("wrong username")
		}
	})
}
