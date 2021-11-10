package refresh_access_token

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/abergmeier/terraform-provider-oauth/internal/hash"
	debuglog "github.com/abergmeier/terraform-provider-oauth/internal/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/oauth2/google"
)

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

	creds, err := google.FindDefaultCredentials(context.TODO())
	if err != nil {
		return nil, err
	}

	c := &defaultCredentials{}
	err = json.Unmarshal(creds.JSON, &c)
	if err != nil {
		return nil, err
	}

	debuglog.DebugLogJSON(creds.JSON)

	return c, nil
}

func read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	iid := d.Get("client_id")
	clientId := iid.(string)
	isecret := d.Get("client_secret")
	clientSecret := isecret.(string)
	rt := d.Get("refresh_token")
	refreshToken := rt.(string)

	if clientId == "" {
		c, err := readDefaultCredentials()
		if err != nil {
			return diag.FromErr(err)
		}
		clientId = c.ClientId

		if clientSecret == "" {
			clientSecret = c.ClientSecret
		}

		if refreshToken == "" {
			refreshToken = c.RefreshToken
		}
	}

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

	hash := hash.BuildHash(clientId, clientSecret, refreshToken, tokenUrl)
	d.SetId(fmt.Sprintf("%x", hash))

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
