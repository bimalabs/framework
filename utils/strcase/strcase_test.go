package strcase

import "testing"

func TestToSnake(t *testing.T) {
	cases := [][2]string{
		{"", ""},
		{" \t\n", ""},
		{"firstName", "first_name"},
		{"FirstName", "first_name"},
		{"already_snake", "already_snake"},
		{" JSONData ", "json_data"},
		{"userID", "user_id"},
		{"HTTPServerURL", "http_server_url"},
		{"AAAbbb", "aa_abbb"},
		{"numbers2and55with000", "numbers_2_and_55_with_000"},
		{"1A2", "1_a_2"},
		{"AB1 AB2", "ab_1_ab_2"},
		{"AnyKind of_string.test-case", "any_kind_of_string_test_case"},
		{"a__b--c..d  e", "a__b__c__d__e"},
		{"a\tb\nc", "a\tb\nc"},
		{"Hello世界", "hello世界"},
		{"ÉcoleName", "École_name"},
		{"\u2003TestCase\u2003", "test_case"},
		{"A/B@C", "a/b@c"},
		{"\xffA\x00B", "\xffa\x00b"},
	}
	for _, test := range cases {
		t.Run(test[0], func(t *testing.T) {
			if got := ToSnake(test[0]); got != test[1] {
				t.Fatalf("ToSnake(%q) = %q, want %q", test[0], got, test[1])
			}
		})
	}
}

func TestToDelimited(t *testing.T) {
	for _, delimiter := range []byte{'_', '-', '.', '/', 0, 0xff} {
		want := "json" + string([]byte{delimiter}) + "data" + string([]byte{delimiter}) + "2" + string([]byte{delimiter}) + "value"
		if got := ToDelimited("JSONData2.value", delimiter); got != want {
			t.Errorf("delimiter %x: got %q, want %q", delimiter, got, want)
		}

		if got := ToDelimited(" \t", delimiter); got != "" {
			t.Errorf("empty input: got %q", got)
		}
	}
}
