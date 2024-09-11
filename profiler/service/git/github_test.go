package git

import (
	"testing"
)

func TestFetchUserDetails(t *testing.T) {
	t.Run("Check user details", func(t *testing.T) {
		gs := NewGitService()
		user, err := gs.FetchUserDetails("xhermitx")
		t.Log(user)
		if err != nil {
			t.Errorf("error occured: %v", err)
		}
	})
}
