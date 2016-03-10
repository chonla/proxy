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
}
