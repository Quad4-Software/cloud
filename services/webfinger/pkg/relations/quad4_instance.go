package relations

import (
	"context"
	"net/url"
	"regexp"
	"strings"
	"text/template"

	"github.com/Quad4-Software/cloud/pkg/oidc"
	"github.com/Quad4-Software/cloud/services/webfinger/pkg/config"
	"github.com/Quad4-Software/cloud/services/webfinger/pkg/service/v0"
	"github.com/Quad4-Software/cloud/services/webfinger/pkg/webfinger"
)

const (
	Quad4CloudInstanceRel = "http://webfinger.quad4/rel/server-instance"
)

type compiledInstance struct {
	config.Instance
	compiledRegex *regexp.Regexp
	hrefTemplate  *template.Template
}

type quad4CloudInstance struct {
	instances    []compiledInstance
	quad4CloudURL string
	instanceHost string
}

// Quad4CloudInstance adds one or more Quad4 Cloud instance relations
func Quad4CloudInstance(instances []config.Instance, quad4CloudURL string) (service.RelationProvider, error) {
	compiledInstances := make([]compiledInstance, 0, len(instances))
	var err error
	for _, instance := range instances {
		compiled := compiledInstance{Instance: instance}
		compiled.compiledRegex, err = regexp.Compile(instance.Regex)
		if err != nil {
			return nil, err
		}
		compiled.hrefTemplate, err = template.New(instance.Claim + ":" + instance.Regex + ":" + instance.Href).Parse(instance.Href)
		if err != nil {
			return nil, err
		}
		compiledInstances = append(compiledInstances, compiled)
	}

	u, err := url.Parse(quad4CloudURL)
	if err != nil {
		return nil, err
	}
	return &quad4CloudInstance{
		instances:    compiledInstances,
		quad4CloudURL: quad4CloudURL,
		instanceHost: u.Host + u.Path,
	}, nil
}

func (l *quad4CloudInstance) Add(ctx context.Context, _ string, jrd *webfinger.JSONResourceDescriptor) {
	if jrd == nil {
		jrd = &webfinger.JSONResourceDescriptor{}
	}
	if claims := oidc.FromContext(ctx); claims != nil {
		if value, ok := claims[oidc.PreferredUsername].(string); ok {
			jrd.Subject = "acct:" + value + "@" + l.instanceHost
		} else if value, ok := claims[oidc.Email].(string); ok {
			jrd.Subject = "mailto:" + value
		}
		// allow referencing OC_URL in the template
		claims["OC_URL"] = l.quad4CloudURL
		for _, instance := range l.instances {
			if value, ok := claims[instance.Claim].(string); ok && instance.compiledRegex.MatchString(value) {
				var tmplWriter strings.Builder
				instance.hrefTemplate.Execute(&tmplWriter, claims)
				jrd.Links = append(jrd.Links, webfinger.Link{
					Rel:    Quad4CloudInstanceRel,
					Href:   tmplWriter.String(),
					Titles: instance.Titles,
				})
				if instance.Break {
					break
				}
			}
		}
	}
}
