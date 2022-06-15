package send

import (
	"testing"
)

// Success and failure markers.
const (
	success = "\u2713"
	failed  = "\u2717"
)

// TestEmail send a test email.
func TestEmail(t *testing.T) {
	t.Log("Given the need to be able to authenticate and authorize access.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen sening an email.", testID)
		{
			if err := Email("example@example.com", "example@example.com", "Test email Serves", "This is a test email"); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to send an email:%v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to send an email.", success, testID)
		}
	}
}
