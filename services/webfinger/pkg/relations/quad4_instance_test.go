package relations

import (
	"context"
	"testing"

	"github.com/Quad4-Software/cloud/pkg/oidc"
	"github.com/Quad4-Software/cloud/services/webfinger/pkg/config"
	"github.com/Quad4-Software/cloud/services/webfinger/pkg/webfinger"
)

func TestQuad4CloudInstanceErr(t *testing.T) {
	_, err := Quad4CloudInstance([]config.Instance{}, "http://\n\rinvalid")
	if err == nil {
		t.Errorf("provider did not err on invalid url: %v", err)
	}
	_, err = Quad4CloudInstance([]config.Instance{{Regex: "("}}, "http://quad4.tld")
	if err == nil {
		t.Errorf("provider did not err on invalid regex: %v", err)
	}
	_, err = Quad4CloudInstance([]config.Instance{{Href: "{{invalid}}ee"}}, "http://quad4.tld")
	if err == nil {
		t.Errorf("provider did not err on invalid href template: %v", err)
	}
}

func TestQuad4CloudInstanceAddLink(t *testing.T) {
	provider, err := Quad4CloudInstance([]config.Instance{{
		Claim: "customclaim",
		Regex: ".+@.+\\..+",
		Href:  "https://{{.otherclaim}}.domain.tld",
		Titles: map[string]string{
			"foo": "bar",
		},
		Break: true,
	}}, "http://quad4.tld")
	if err != nil {
		t.Error(err)
	}

	ctx := context.Background()
	ctx = oidc.NewContext(ctx, map[string]any{
		"customclaim": "some@fizz.buzz",
		"otherclaim":  "someone",
	})
	jrd := webfinger.JSONResourceDescriptor{}
	provider.Add(ctx, "", &jrd)

	if len(jrd.Links) != 1 {
		t.Errorf("provider returned wrong number of links: %v, expected 1", len(jrd.Links))
	}
	if jrd.Links[0].Href != "https://someone.domain.tld" {
		t.Errorf("provider returned wrong issuer link href: %v, expected %v", jrd.Links[0].Href, "https://someone.domain.tld")
	}
	if jrd.Links[0].Rel != Quad4CloudInstanceRel {
		t.Errorf("provider returned quad4 server instance rel: %v, expected %v", jrd.Links[0].Rel, Quad4CloudInstanceRel)
	}
	if len(jrd.Links[0].Titles) != 1 {
		t.Errorf("provider returned wrong number of titles: %v, expected 1", len(jrd.Links[0].Titles))
	}
	if jrd.Links[0].Titles["foo"] != "bar" {
		t.Errorf("provider returned wrong title: %v, expected bar", len(jrd.Links[0].Titles["foo"]))
	}

}
