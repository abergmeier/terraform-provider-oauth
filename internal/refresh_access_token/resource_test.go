package refresh_access_token

import (
	"strings"
	"testing"
)

func TestSetDataFromReader(t *testing.T) {

	r := strings.NewReader(`{
"access_token": "1/fFAGRNJru1FTz70BzhT3Zg",
"expires_in": 3920,
"scope": "https://www.googleapis.com/auth/drive.metadata.readonly",
"token_type": "Bearer"
}`)
	d := Resource().TestResourceData()
	diags := setDataFromReader(r, d)
	if diags != nil {
		t.Error("Unexpected problems:", diags)
	}

	if d.Get("access_token").(string) != "1/fFAGRNJru1FTz70BzhT3Zg" {
		t.Error("Invalid access_token")
	}

	if d.Get("id_token").(string) != "" {
		t.Error("Invalid id_token")
	}

	if d.Get("scope").(string) != "https://www.googleapis.com/auth/drive.metadata.readonly" {
		t.Error("Invalid scope")
	}

	if d.Get("token_type").(string) != "Bearer" {
		t.Error("Invalid token_type")
	}
}
