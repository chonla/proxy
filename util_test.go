package main

import "testing"

func TestFoundInHostList(t *testing.T) {
	var testcases = []struct {
		name     string
		expected bool
		list     string
		search   string
	}{
		{"match", true, "google.com,yahoo.com", "google.com"},
		{"match", true, "google.com,yahoo.com", "yahoo.com"},
		{"not match", false, "google.com,yahoo.com", "bing.com"},
		{"not match", false, "", "bing.com"},
	}

	for _, testcase := range testcases {
		if inHostList(testcase.list, testcase.search) != testcase.expected {
			t.Error("fail case ", testcase.name)
		}
	}
}

func TestChangeHostToHttps(t *testing.T) {
	result := changeHostToHttps("http://www.google.com:9000")
	if result != "https://www.google.com:9000" {
		t.Error("expect url https but got ", result)
	}
}

func TestGetValueByKeyName(t *testing.T) {
	data := "<urn:SUBR_NUMB>66897893394</urn:SUBR_NUMB>"

	result := getValueByKey("urn:SUBR_NUMB", data)

	if result != "66897893394" {
		t.Error("expect 66897893394 but got", result)
	}
}

func TestGetValueByKeyNameBlank(t *testing.T) {
	data := "<urn:SUBR_NUMB></urn:SUBR_NUMB>"

	result := getValueByKey("urn:SUBR_NUMB", data)

	if result != "" {
		t.Error("expect <blank> but got", result)
	}

	result = getValueByKey("urn:CUST_NUMB", data)

	if result != "" {
		t.Error("expect <blank> but got", result)
	}

}

func TestGetValueFromRAW(t *testing.T) {
	data := `<bbm:QueryEstimatedChargeRequest>
           <!--Optional:-->
           <ccin:SubscriberNo>999850151</ccin:SubscriberNo>
           <!--Optional:-->
           <ccin:AccountCode>xxxxxxx</ccin:AccountCode>
           <ccin:LanguageType>1</ccin:LanguageType>
           <ccin:BillCycleMonth>201501</ccin:BillCycleMonth>
        </bbm:QueryEstimatedChargeRequest>`

	result := getValueByKey("ccin:AccountCode", data)

	if result != "xxxxxxx" {
		t.Error("expect xxxxxxx but got", result)
	}
}

func TestGetConditionField(t *testing.T) {
	list := Condition{}
	list["AAA"] = "AAA"
	list["BBB"] = "BBB"

	result := getConditionField("AAA", list)

	if result != "AAA" {
		t.Error("expect AAA but got", result)
	}

	result = getConditionField("XXX", list)

	if result != "" {
		t.Error("expect <blank> but got", result)
	}
}

func TestGetConditionValue(t *testing.T) {
	result := getConditionValue("A", "<A>xxx</A>")
	if result != "xxx" {
		t.Error("expect xxx but got", result)
	}

	result = getConditionValue("A,C", "<A>xxx</A><B>yyy</B><C>zzz</C>")
	if result != "xxx|zzz" {
		t.Error("expect xxx|zzz but got", result)
	}

	result = getConditionValue("A,M", "<A>xxx</A><B>yyy</B><C>zzz</C>")
	if result != "xxx" {
		t.Error("expect xxx but got", result)
	}

}

func TestFoundInIncludeList(t *testing.T) {
	condition := make(map[string]string)
	condition["www.google.com/starwars"] = "AAA"
	condition["dtac.co.th/prepaid"] = "BBB"
	var testcases = []struct {
		name     string
		expected bool
		search   Recoder
	}{
		{"match1", true, Recoder{Request: Inbound{Host: "www.google.com", Path: "/starwars"}}},
		{"match2", true, Recoder{Request: Inbound{Host: "dtac.co.th", Path: "/prepaid"}}},
		{"not match1", false, Recoder{Request: Inbound{Host: "bing.com", Path: "/reward"}}},
		{"not match2", false, Recoder{Request: Inbound{Host: "yahoo.com", Path: "/mail"}}},
	}

	original := arg.IncludeList
	arg.IncludeList = condition
	for _, testcase := range testcases {
		if foundIncludeList(testcase.search) != testcase.expected {
			t.Error("fail case ", testcase.name)
		}
	}
	arg.IncludeList = original
}

func TestWildcardIncludeList(t *testing.T) {
	condition := make(map[string]string)
	condition["www.google.com/star*"] = "A"
	condition["yahoo.com/star"] = "B"
	var testcases = []struct {
		name     string
		expected bool
		search   string
	}{
		{"match1", true, "www.google.com/starwars"},
		{"match2", true, "www.google.com/star"},
		{"not match1", false, "google.com/star"},
		{"not match2", false, "th.yahoo.com/startrek"},
	}

	original := arg.IncludeList
	arg.IncludeList = condition
	for _, testcase := range testcases {
		if wildcardMatch(testcase.search) != testcase.expected {
			t.Error("fail case ", testcase.name)
		}
	}
	arg.IncludeList = original
}

func TestIsRegX(t *testing.T) {
	if !hasAsteriskCharactor("google.com/star*") {
		t.Error("expect has AsteriskCharactor")
	}

	if !isRegularExpression("google.com/star*") {
		t.Error("expect equal RegularExpression")
	}

	if isRegularExpression("google.com/star") {
		t.Error("expect not equal RegularExpression")
	}
}
