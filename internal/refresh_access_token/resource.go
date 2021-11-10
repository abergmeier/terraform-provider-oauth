package refresh_access_token

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Client ID",
				Sensitive:   true,
				DefaultFunc: readDefaultGcloudClientId,
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Client Secret",
				Sensitive:   true,
				DefaultFunc: readDefaultGcloudClientSecret,
			},
			"refresh_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Refresh Token",
				Sensitive:   true,
				DefaultFunc: readDefaultGcloudRefreshToken,
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
			"id_token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"token_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		ReadContext: read,
	}
}

type defaultCredentials struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
}
type refreshResponse struct {
	AccessToken string `json:"access_token"`
	IdToken     string `json:"id_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func readDefaultCredentials() (*defaultCredentials, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(fmt.Sprintf("%s/.config/gcloud/application_default_credentials.json", homeDir))
	if err != nil {
		return nil, err
	}

	c := &defaultCredentials{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func readDefaultGcloudClientId() (interface{}, error) {
	c, err := readDefaultCredentials()
	if err != nil {
		return nil, err
	}
	return c.ClientId, nil
}

func readDefaultGcloudClientSecret() (interface{}, error) {
	c, err := readDefaultCredentials()
	if err != nil {
		return nil, err
	}
	return c.ClientSecret, nil
}

func readDefaultGcloudRefreshToken() (interface{}, error) {
	c, err := readDefaultCredentials()
	if err != nil {
		return nil, err
	}
	return c.RefreshToken, nil
}

func read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	iid := d.Get("client_id")
	clientId := iid.(string)
	isecret := d.Get("client_secret")
	clientSecret := isecret.(string)
	rt := d.Get("refresh_token")
	refreshToken := rt.(string)
	tu := d.Get("token_url")
	tokenUrl := tu.(string)

	p := fmt.Sprintf("client_id=%s&client_secret=%s&refresh_token=%s&grant_type=refresh_token", clientId, clientSecret, refreshToken)
	r, err := http.Post(
		tokenUrl,
		"application/x-www-form-urlencoded",
		strings.NewReader(p),
	)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	return setDataFromResponse(r, d)
}

func debugLogResponse(s []byte) {
	if len(s) > 1024 {
		log.Printf(`[DEBUG] OAuth response details (cropped):
---[ RESPONSE ]--------------------------------------
%s
...
-----------------------------------------------------
`, string(s[:1024]))
	} else {
		log.Printf(`[DEBUG] OAuth response details:
---[ RESPONSE ]--------------------------------------
%s
-----------------------------------------------------
`, string(s))
	}
}

func setDataFromResponse(r *http.Response, d *schema.ResourceData) diag.Diagnostics {

	rb, err := io.ReadAll(r.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	debugLogResponse(rb)

	if r.StatusCode < 200 || r.StatusCode > 299 {
		return diag.Errorf("Responded with code %d: %s\n", r.StatusCode, http.StatusText(r.StatusCode))
	}

	return setDataFromJSON(rb, d)
}

func setDataFromJSON(s []byte, d *schema.ResourceData) diag.Diagnostics {

	c := &refreshResponse{}
	err := json.Unmarshal(s, c)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("access_token", c.AccessToken)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("id_token", c.IdToken)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("scope", c.Scope)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("token_type", c.TokenType)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
