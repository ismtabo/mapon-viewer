package steps

import (
	"context"
	"time"

	"github.com/Telefonica/golium"
	"github.com/Telefonica/golium/steps/http"
	"github.com/cucumber/godog"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// ServerSteps add golium steps for ATs.
type ServerSteps struct{}

// InitializeSteps add steps additional godog steps.
func (s ServerSteps) InitializeSteps(ctx context.Context, scenCtx *godog.ScenarioContext) context.Context {
	// Retrieve HTTP session
	session := http.GetSession(ctx)
	// Initialize the steps
	scenCtx.Step(`^the HTTP response should not be empty$`, func() error {
		if err := session.ValidateResponseBodyEmpty(ctx); err == nil {
			return errors.New("error: HTTP response is empty")
		}
		return nil
	})
	scenCtx.Step(`^the HTTP response should containt the text$`, func(t *godog.DocString) error {
		body := golium.ValueAsString(ctx, t.Content)
		return ValidateResponseTextContent(ctx, body)
	})
	scenCtx.Step(`^the HTTP response should contain the( expired)? cookie "([^"]+)"$`, func(expires, name string) error {
		checkExpired := expires != ""
		name = golium.ValueAsString(ctx, name)
		cookies := session.Response.Response.Cookies()
		for _, cookie := range cookies {
			if cookie.Name == name {
				if checkExpired && cookie.Expires.After(time.Now()) {
					return errors.Errorf("failed validating http response cookie. cookie with name '%s' has not expired", name)
				}
				if !checkExpired && cookie.Expires.Before(time.Now()) {
					return errors.Errorf("failed validating http response cookie. cookie with name '%s' has expired", name)
				}
				return nil
			}
		}
		return errors.Errorf("failed validating http response cookie. cookie not found with name '%s'", name)
	})
	scenCtx.Step(`^I keep the HTTP cookies$`, func() error {
		if session.Request.Headers == nil {
			session.Request.Headers = make(map[string][]string)
		}
		cookiesHeader, found := session.Request.Headers["Cookie"]
		if !found {
			cookiesHeader = make([]string, 0)
		}
		cookies := session.Response.Response.Cookies()
		for _, cookie := range cookies {
			cookiesHeader = append(cookiesHeader, cookie.String())
		}
		session.Request.Headers["Cookie"] = cookiesHeader
		return nil
	})
	scenCtx.Step(`^I generate a salt hash for password "([^"]+)" and store in context "([^"]+)"$`, func(pwd, key string) error {
		pwd = golium.ValueAsString(ctx, pwd)
		key = golium.ValueAsString(ctx, key)
		password, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
		if err != nil {
			return err
		}
		golium.GetContext(ctx).Put(key, string(password))
		return nil
	})
	return ctx
}

// ValidateResponseTextContent validates Http session response has the expected text body content.
func ValidateResponseTextContent(ctx context.Context, expected string) error {
	session := http.GetSession(ctx)
	actual := string(session.Response.ResponseBody)
	if expected != actual {
		return errors.Errorf("failed validating response text body: expected '%s' actual '%s'", expected, actual)
	}
	return nil
}
