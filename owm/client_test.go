package owm

import "testing"

func TestBuildUrl(t *testing.T) {
	expected := "http://api.openweathermap.org/data/2.5/weather?id=123&appid=foobar&lang=de&units=metric"
	got := buildUrl(123, "foobar", "de")
	if expected != got {
		t.Errorf("Expected: %s, Got: %s", expected, got)
	}
}
