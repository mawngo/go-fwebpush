package fastunsafeurl

import "testing"

func TestParseSchemeHost(t *testing.T) {
	cases := [][]any{
		{"https://example.com", "https://example.com", true},
		{"https://example.com/test123", "https://example.com", true},
		{"abc://example.com/test123#frag?abc=1", "abc://example.com", true},
		{"://", "", false},
		{"example.com", "", false},
		{"a://b", "a://b", true},
		{"http://user:account@example.com", "http://user:account@example.com", true},
		{"http://user:account@example.com:8080", "http://user:account@example.com:8080", true},
		{"https://example.com:8080", "https://example.com:8080", true},
		{"https://updates.push.services.mozilla.com/wpush/v2/gAAAAA", "https://updates.push.services.mozilla.com", true},
		{"https://fcm.googleapis.com/fcm/send/eKAWKNUIYFw:APA91bHkYaziMvso61arnA20A8j83Mv7uv8ud", "https://fcm.googleapis.com", true},
	}
	for i := range cases {
		endpoint := cases[i][0].(string)
		expected := cases[i][1].(string)
		valid := cases[i][2].(bool)
		actual, _, err := ParseSchemeHost(endpoint)

		if !valid {
			if actual != "" {
				t.Fatal("Expected empty audience from", endpoint, ":", actual)
			}
			if err == nil {
				t.Fatal("Expected error from", endpoint)
			}
			continue
		}
		if expected != actual {
			t.Fatal("Expected", expected, "from", endpoint, "got", actual)
		}
		if err != nil {
			t.Fatal("Expected no error from", endpoint, "got", err)
		}
	}
}

func TestParseHost(t *testing.T) {
	cases := [][]any{
		{"https://example.com", "example.com", true},
		{"https://example.com/test123", "example.com", true},
		{"abc://example.com/test123#frag?abc=1", "example.com", true},
		{"://", "", false},
		{"example.com", "", false},
		{"a://b", "b", true},
		{"http://user:account@example.com", "user:account@example.com", true},
		{"http://user:account@example.com:8080", "user:account@example.com:8080", true},
		{"https://example.com:8080", "example.com:8080", true},
		{"https://updates.push.services.mozilla.com/wpush/v2/gAAAAA", "updates.push.services.mozilla.com", true},
		{"https://fcm.googleapis.com/fcm/send/eKAWKNUIYFw:APA91bHkYaziMvso61arnA20A8j83Mv7uv8ud", "fcm.googleapis.com", true},
	}
	for i := range cases {
		endpoint := cases[i][0].(string)
		expected := cases[i][1].(string)
		valid := cases[i][2].(bool)
		actual, err := ParseHost(endpoint)

		if !valid {
			if actual != "" {
				t.Fatal("Expected empty audience from", endpoint, ":", actual)
			}
			if err == nil {
				t.Fatal("Expected error from", endpoint)
			}
			continue
		}
		if expected != actual {
			t.Fatal("Expected", expected, "from", endpoint, "got", actual)
		}
		if err != nil {
			t.Fatal("Expected no error from", endpoint, "got", err)
		}
	}
}
