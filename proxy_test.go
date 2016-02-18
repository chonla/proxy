package main

import "testing"

func TestGetValueByKeyName(t *testing.T) {
	data := "<urn:SUBR_NUMB>66897893394</urn:SUBR_NUMB>"

	result := getValueByKey("SUBR_NUMB", data)

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

	result := getValueByKey("AccountCode", data)

	if result != "xxxxxxx" {
		t.Error("expect xxxxxxx but got", result)
	}
}
