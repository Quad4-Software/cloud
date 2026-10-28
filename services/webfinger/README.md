# Webfinger

The webfinger service provides an RFC7033 WebFinger lookup of Quad4 Cloud resources, relevant for a given user account at the /.well-known/webfinger enpoint.

1. An [OpenID Connect Discovery](#openid-connect-discovery) for the IdP, based on the Quad4 Cloud URL.
2. An [Authenticated Instance Discovery](#authenticated-instance-discovery), based on the user account.

These two request are only needed for discovery.

## OpenID Connect Discovery

Clients can make an unauthenticated `GET https://drive.quad4.test/.well-known/webfinger?resource=https%3A%2F%2Fdrive.quad4.test` request to discover the OpenID Connect Issuer in the `http://openid.net/specs/connect/1.0/issuer` relation:

```json
{
    "subject": "https://drive.quad4.test",
    "links": [
        {
            "rel": "http://openid.net/specs/connect/1.0/issuer",
            "href": "https://sso.example.org/cas/oidc/"
        }
    ]
}
```

Here, the `resource` takes the instance domain URI, but an `acct:` URI works as well.

## Authenticated Instance Discovery

When using OpenID connect to authenticate requests, clients can look up the quad4 instances a user has access to.

*   Authentication is necessary to prevent leaking information about existing users.
*   Basic auth is not supported.

The default configuration will simply return the `OC_URL` and direct clients to that domain:

```json
{
    "subject": "acct:alan@drive.quad4.test",
    "links": [
        {
            "rel": "http://openid.net/specs/connect/1.0/issuer",
            "href": "https://sso.example.org/cas/oidc/"
        },
        {
            "rel": "http://webfinger.quad4/rel/server-instance",
            "href": "https://abc.drive.example.org",
            "titles": {
                "en": "Quad4 Cloud Instance"
            }
        }
    ]
}
```

## Configure Different Instances Based on OpenidConnect UserInfo Claims

A more complex example for configuring different instances could look like this:

```yaml
webfinger:
  instances:
  -  claim: email
     regex: alan@example\.org
     href: "https://{{.preferred_username}}.cloud.quad4.test"
     title:
       "en": "Quad4 Cloud Instance for Alan"
       "de": "Quad4 Cloud Instanz für Alan"
     break: true
  -  claim: "email"
     regex: mary@example\.org
     href: "https://{{.preferred_username}}.cloud.quad4.test"
     title:
       "en": "Quad4 Cloud Instance for Mary"
       "de": "Quad4 Cloud Instanz für Mary"
     break: false
  -  claim: "email"
     regex: .+@example\.org
     href: "https://example-org.cloud.quad4.test"
     title:
       "en": "Quad4 Cloud Instance for example.org"
       "de": "Quad4 Cloud Instanz für example.org"
     break: true
  -  claim: "email"
     regex: .+@example\.com
     href: "https://example-com.cloud.quad4.test"
     title:
       "en": "Quad4 Cloud Instance for example.com"
       "de": "Quad4 Cloud Instanz für example.com"
     break: true
  -  claim: "email"
     regex: .+@.+\..+
     href: "https://cloud.quad4.test"
     title:
       "en": "Quad4 Cloud Instance"
       "de": "Quad4 Cloud Instanz"
     break: true
```

Now, an authenticated webfinger request for `acct:me@example.org` (when logged in as mary) would return two instances, based on her `email` claim, the regex matches and break flags:

```json
{
    "subject": "acct:mary@example.org",
    "links": [
        {
            "rel": "http://openid.net/specs/connect/1.0/issuer",
            "href": "https://sso.example.org/cas/oidc/"
        },
        {
            "rel": "http://webfinger.quad4/rel/server-instance",
            "href": "https://mary.cloud.quad4.test",
            "titles": {
                "en": "Quad4 Cloud Instance for Mary",
                "de": "Quad4 Cloud Instanz für Mary"
            }
        },
        {
            "rel": "http://webfinger.quad4/rel/server-instance",
            "href": "https://xyz.drive.example.org",
            "titles": {
                "en": "Quad4 Cloud Instance for example.org",
                "de": "Quad4 Cloud Instanz für example.org"
            }
        }
    ]
}
```
