package jwt

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/shulganew/GophKeeper/internal/app/config"
)

// JWSValidator is used to validate JWS payloads and return a JWT if they're
// valid
type JWSValidator interface {
	ValidateJWS(jws string) (jwt.Token, error)
}

// GetJWSFromRequest extracts a JWS string from an Authorization: Bearer <jws> header
func GetJWSFromRequest(req *http.Request) (string, error) {
	authHdr := req.Header.Get("Authorization")
	// Check for the Authorization header.
	if authHdr == "" {
		return "", errors.New("authorization header is missing")
	}
	// We expect a header value of the form "Bearer <token>", with 1 space after
	// Bearer, per spec.
	prefix := "Bearer "
	if !strings.HasPrefix(authHdr, prefix) {
		return "", errors.New("authorization header is malformed")
	}
	return strings.TrimPrefix(authHdr, prefix), nil
}

// Get user uuid.
func GetUserID(v JWSValidator, req *http.Request) (string, error) {
	jws, err := GetJWSFromRequest(req)
	if err != nil {
		return "", fmt.Errorf("getting jws: %w", err)
	}

	// if the JWS is valid, we have a JWT, which will contain a bunch of claims.
	token, err := v.ValidateJWS(jws)
	if err != nil {
		return "", fmt.Errorf("validating JWS: %w", err)
	}
	userIDInt, isExisted := token.Get(UserIdentification)
	if !isExisted {
		return "", fmt.Errorf("userID not found: %w", err)
	}
	iserID, ok := userIDInt.(string)
	if !ok {
		return "", fmt.Errorf("userID can't cast to string: %w", err)
	}
	return iserID, nil

}

func NewAuthenticator(v JWSValidator) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		return Authenticate(v, ctx, input)
	}
}

// Authenticate uses the specified validator to ensure a JWT is valid, then makes
// sure that the claims provided by the JWT match the scopes as required in the API.
func Authenticate(v JWSValidator, ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	// Our security scheme is named BearerAuth, ensure this is the case
	if input.SecuritySchemeName != "BearerAuth" {
		return fmt.Errorf("security scheme %s != 'BearerAuth'", input.SecuritySchemeName)
	}

	// Now, we need to get the JWS from the request, to match the request expectations
	// against request contents.
	jws, err := GetJWSFromRequest(input.RequestValidationInput.Request)
	if err != nil {
		return fmt.Errorf("getting jws: %w", err)
	}

	// if the JWS is valid, we have a JWT, which will contain a bunch of claims.
	token, err := v.ValidateJWS(jws)
	if err != nil {
		return fmt.Errorf("validating JWS: %w", err)
	}

	// We've got a valid token now, and we can look into its claims to see whether
	// they match. Every single scope must be present in the claims.
	err = CheckTokenClaims(input.Scopes, token)

	if err != nil {
		return fmt.Errorf("token claims don't match: %w", err)
	}

	return nil
}

// GetClaimsFromToken returns a list of claims from the token. We store these
// as a list under the "perms" claim, short for permissions, to keep the token
// shorter.
func GetClaimsFromToken(t jwt.Token) ([]string, error) {
	rawPerms, found := t.Get(PermissionsClaim)
	if !found {
		// If the perms aren't found, it means that the token has none, but it has
		// passed signature validation by now, so it's a valid token, so we return
		// the empty list.
		return make([]string, 0), nil
	}

	// rawPerms will be an untyped JSON list, so we need to convert it to
	// a string list.
	rawList, ok := rawPerms.([]interface{})
	if !ok {
		return nil, fmt.Errorf("'%s' claim is unexpected type'", PermissionsClaim)
	}

	claims := make([]string, len(rawList))

	for i, rawClaim := range rawList {
		var ok bool
		claims[i], ok = rawClaim.(string)
		if !ok {
			return nil, fmt.Errorf("%s[%d] is not a string", PermissionsClaim, i)
		}
	}
	return claims, nil
}

func CheckTokenClaims(expectedClaims []string, t jwt.Token) error {
	claims, err := GetClaimsFromToken(t)
	if err != nil {
		return fmt.Errorf("getting claims from token: %w", err)
	}
	// Put the claims into a map, for quick access.
	claimsMap := make(map[string]bool, len(claims))
	for _, c := range claims {
		claimsMap[c] = true
	}

	for _, e := range expectedClaims {
		if !claimsMap[e] {
			return errors.New("provided claims do not match expected scopes")
		}
	}
	return nil
}
func GetPrivateKey(conf config.Config) (key []byte, err error) {
	file, err := os.Open(conf.PathJWT)
	if err != nil {
		return nil, fmt.Errorf("couldn't open file from path: %w", err)
	}
	key, err = io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("couldn't read from file: %w", err)
	}
	return key, nil
}
