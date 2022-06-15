package feed

import (
	"testing"
)

// Success and failure markers.
const (
	success = "\u2713"
	failed  = "\u2717"
)

// TestGeolocation test the ability to feed system with geolocation data.
func TestGeolocation(t *testing.T) {
	t.Log("Given the need to be able to authenticate and authorize access.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a geolocation.", testID)
		{

			geo, err := GetGeolocation("178.130.111.140")
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to get the geolocation for the ip address : %s.", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to get the geolocation for the ip address : %v.", success, testID, geo)
		}
	}
}
