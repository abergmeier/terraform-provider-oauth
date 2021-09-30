package refresh_access_token

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Client ID",
				Sensitive:   true,
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Client Secret",
				Sensitive:   true,
			},
			"refresh_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Refresh Token",
				Sensitive:   true,
				DefaultFunc: gcloudRefreshToken,
			},
			"token_url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https://oauth2.googleapis.com/token",
			},
			"access_token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"token_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		ReadContext: read,
	}
}

type applicationDefaultCredentials struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func gcloudRefreshToken() (interface{}, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(fmt.Sprintf("%s/.config/gcloud/application_default_credentials.json", homeDir))
	if err != nil {
		return nil, err
	}

	c := &applicationDefaultCredentials{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}

	return c.AccessToken, nil
}

func read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	iid := d.Get("client_id")
	clientId := iid.(string)
	isecret := d.Get("client_secret")
	clientSecret := isecret.(string)
	rt := d.Get("refresh_token")
	refreshToken := rt.(string)
	tu := d.Get("token_url")
	tokenUrl := tu.(string)

	r, err := http.PostForm(tokenUrl, url.Values{
		"client_id":     {clientId},
		"client_secret": {clientSecret},
		"refresh_token": {refreshToken},
		"grant_type":    {"refresh_token"},
	})
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	defer r.Body.Close()

	rb, err := io.ReadAll(r.Body)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	c := &applicationDefaultCredentials{}
	err = json.Unmarshal(rb, c)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	err = d.Set("access_token", c.AccessToken)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	err = d.Set("scope", c.Scope)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	err = d.Set("token_type", c.TokenType)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}
