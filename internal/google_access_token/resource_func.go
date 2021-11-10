package google_access_token

import (
	"context"
	"fmt"

	"github.com/abergmeier/terraform-provider-oauth/internal/hash"
	"github.com/abergmeier/terraform-provider-oauth/internal/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	scopes := d.Get("scopes").([]string)

	creds, err := google.FindDefaultCredentials(context.TODO(), scopes...)
	if err != nil {
		return diag.FromErr(err)
	}

	log.DebugLogJSON(creds.JSON)

	hash := hash.BuildHash(string(creds.JSON))

	d.SetId(fmt.Sprintf("%x", hash))

	return diag.FromErr(setDataFromTokenSource(creds.TokenSource, d))
}

func setDataFromTokenSource(ts oauth2.TokenSource, d *schema.ResourceData) error {

	t, err := ts.Token()
	if err != nil {
		return err
	}

	err = d.Set("access_token", t.AccessToken)
	if err != nil {
		return err
	}

	return d.Set("token_type", t.TokenType)
}
