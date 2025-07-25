package gcorecloud

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// DefaultUserAgent is the default User-Agent string set in the request header.
const DefaultUserAgent = "cloud-api-go-sdk/%s"

var AppVersion = "0.22.0"

// UserAgent represents a User-Agent header.
type UserAgent struct {
	// prepend is the slice of User-Agent strings to prepend to DefaultUserAgent.
	// All the strings to prepend are accumulated and prepended in the Join method.
	prepend []string
}

// Prepend prepends a user-defined string to the default User-Agent string. Users
// may pass in one or more strings to prepend.
func (ua *UserAgent) Prepend(s ...string) {
	ua.prepend = append(s, ua.prepend...)
}

// Join concatenates all the user-defined User-Agent strings with the default
// GCore cloud User-Agent string.
func (ua *UserAgent) Join() string {
	// nolint:gocritic
	uaSlice := append(ua.prepend, fmt.Sprintf(DefaultUserAgent, AppVersion))
	return strings.Join(uaSlice, " ")
}

// ProviderClient stores details that are required to interact with any
// services within a specific provider's API.
//
// Generally, you acquire a ProviderClient by calling the NewClient method in
// the appropriate provider's child package, providing whatever authentication
// credentials are required.
type ProviderClient struct {
	// IdentityBase is the base URL used for a particular provider's identity
	// service - it will be used when issuing authentication requests. It
	// should point to the root resource of the identity service, not a specific
	// identity version.
	IdentityBase string

	// APIURL is the identity endpoint. This may be a specific version
	// of the identity service. If this is the case, this endpoint is used rather
	// than querying versions first.
	IdentityEndpoint string

	// AccessTokenID and RefreshTokenID is the IDs of the most recently issued valid tokens.
	// NOTE: Aside from within a custom ReauthFunc, this field shouldn't be set by an application.
	// To safely read or write this value, call `AccessToken` or `SetAccessToken`
	// call `RefreshToken` or `SetRefreshToken`, respectively
	AccessTokenID  string
	RefreshTokenID string

	// EndpointLocator describes how this provider discovers the endpoints for
	// its constituent services.
	EndpointLocator EndpointLocator

	// HTTPClient allows users to interject arbitrary http, https, or other transit behaviors.
	HTTPClient http.Client

	// UserAgent represents the User-Agent header in the HTTP request.
	UserAgent UserAgent

	// ReauthFunc is the function used to re-authenticate the user if the request
	// fails with a 401 HTTP response code. This a needed because there may be multiple
	// authentication functions for different Identity service versions.
	ReauthFunc func() error

	// Throwaway determines whether if this client is a throw-away client. It's a copy of user's provider client
	// with the token and reauth func zeroed. Such client can be used to perform reauthorization.
	Throwaway bool

	// Context is the context passed to the HTTP request.
	Context context.Context

	// mut is a mutex for the client. It protects read and write access to client attributes such as getting
	// and setting the AccessTokenID.
	mut *sync.RWMutex

	// reauthmut is a mutex for reauthentication it attempts to ensure that only one reauthentication
	// attempt happens at one time.
	reauthmut *reauthlock

	authResult AuthResult

	debug    bool
	APIToken string
	APIBase  string

	// retryGetOn5XX enables GET retries on 5XX errors with the specified number of attempts and a base interval
	// following an exponential backoff with jitter.
	// See EnableGetRetriesOn5XX for enabling this feature.
	retryGetOn5XX             bool
	retryGetOn5XXAttempts     int
	retryGetOn5XXBaseInterval int
}

// reauthlock represents a set of attributes used to help in the reauthentication process.
type reauthlock struct {
	sync.RWMutex
	// This channel is non-nil during reauthentication. It can be used to ask the
	// goroutine doing Reauthenticate() for its result. Look at the implementation
	// of Reauthenticate() for details.
	ongoing chan<- (chan<- error)
}

// AuthenticatedHeaders returns a map of HTTP headers that are common for all
// authenticated service requests. Blocks if Reauthenticate is in progress.
func (client *ProviderClient) AuthenticatedHeaders() (m map[string]string) {
	if client.APIToken != "" {
		return map[string]string{"Authorization": fmt.Sprintf("APIKey %s", client.APIToken)}
	}
	if client.IsThrowaway() {
		return
	}
	if client.reauthmut != nil {
		// If a Reauthenticate is in progress, wait for it to complete.
		client.reauthmut.Lock()
		ongoing := client.reauthmut.ongoing
		client.reauthmut.Unlock()
		if ongoing != nil {
			responseChannel := make(chan error)
			ongoing <- responseChannel
			<-responseChannel
		}
	}
	t := client.AccessToken()
	if t == "" {
		return
	}
	return map[string]string{"Authorization": fmt.Sprintf("Bearer %s", t)}
}

// UseTokenLock creates a mutex that is used to allow safe concurrent access to the auth token.
// If the application's ProviderClient is not used concurrently, this doesn't need to be called.
func (client *ProviderClient) UseTokenLock() {
	client.mut = new(sync.RWMutex)
	client.reauthmut = new(reauthlock)
}

// GetAuthResult returns the result from the request that was used to obtain a
// provider client's token.
//
// The result is nil when authentication has not yet taken place, when the token
// was set manually with SetToken(), or when a ReauthFunc was used that does not
// record the AuthResult.
func (client *ProviderClient) GetAuthResult() AuthResult {
	if client.mut != nil {
		client.mut.RLock()
		defer client.mut.RUnlock()
	}
	return client.authResult
}

// AccessToken safely reads the value of the auth token from the ProviderClient. Applications should
// call this method to access the token instead of the AccessTokenID field
func (client *ProviderClient) AccessToken() string {
	if client.mut != nil {
		client.mut.RLock()
		defer client.mut.RUnlock()
	}
	return client.AccessTokenID
}

// RefreshToken safely reads the value of the auth token from the ProviderClient. Applications should
// call this method to access the token instead of the RefreshTokenID field
func (client *ProviderClient) RefreshToken() string {
	if client.mut != nil {
		client.mut.RLock()
		defer client.mut.RUnlock()
	}
	return client.RefreshTokenID
}

// SetAPIToken safely sets the value of the api token in the ProviderClient
func (client *ProviderClient) SetAPIToken(opt APITokenOptions) error {
	client.APIToken = opt.APIToken
	return nil
}

// SetTokensAndAuthResult safely sets the value of the auth token in the
// ProviderClient and also records the AuthResult that was returned from the
// token creation request. Applications may call this in a custom ReauthFunc.
func (client *ProviderClient) SetTokensAndAuthResult(r AuthResult) error {
	accessTokenID := ""
	refreshTokenID := ""
	var err error
	if r != nil {
		accessTokenID, refreshTokenID, err = r.ExtractTokensPair()
		if err != nil {
			return err
		}
	}

	if client.mut != nil {
		client.mut.Lock()
		defer client.mut.Unlock()
	}
	client.AccessTokenID = accessTokenID
	client.RefreshTokenID = refreshTokenID
	client.authResult = r
	return nil
}

// CopyTokensFrom safely copies the token from another ProviderClient into the
// this one.
func (client *ProviderClient) CopyTokensFrom(other *ProviderClient) {
	if client.mut != nil {
		client.mut.Lock()
		defer client.mut.Unlock()
	}
	if other.mut != nil && other.mut != client.mut {
		other.mut.RLock()
		defer other.mut.RUnlock()
	}
	client.AccessTokenID = other.AccessTokenID
	client.RefreshTokenID = other.RefreshTokenID
	client.authResult = other.authResult
}

// IsThrowaway safely reads the value of the client Throwaway field.
func (client *ProviderClient) IsThrowaway() bool {
	if client.reauthmut != nil {
		client.reauthmut.RLock()
		defer client.reauthmut.RUnlock()
	}
	return client.Throwaway
}

// SetThrowaway safely sets the value of the client Throwaway field.
func (client *ProviderClient) SetThrowaway(v bool) {
	if client.reauthmut != nil {
		client.reauthmut.Lock()
		defer client.reauthmut.Unlock()
	}
	client.Throwaway = v
}

// Reauthenticate calls client.ReauthFunc in a thread-safe way. If this is
// called because of a 401 response, the caller may pass the previous token. In
// this case, the reauthentication can be skipped if another thread has already
// reauthenticated in the meantime. If no previous token is known, an empty
// string should be passed instead to force unconditional reauthentication.
func (client *ProviderClient) Reauthenticate(previousToken string) error {
	if client.ReauthFunc == nil {
		return nil
	}

	if client.reauthmut == nil {
		return client.ReauthFunc()
	}

	messages := make(chan (chan<- error))

	// Check if a Reauthenticate is in progress, or start one if not.
	client.reauthmut.Lock()
	ongoing := client.reauthmut.ongoing
	if ongoing == nil {
		client.reauthmut.ongoing = messages
	}
	client.reauthmut.Unlock()

	// If Reauthenticate is running elsewhere, wait for its result.
	if ongoing != nil {
		responseChannel := make(chan error)
		ongoing <- responseChannel
		return <-responseChannel
	}

	// Perform the actual reauthentication.
	var err error
	if previousToken == "" || client.AccessTokenID == previousToken {
		err = client.ReauthFunc()
	} else {
		err = nil
	}

	// Mark Reauthenticate as finished.
	client.reauthmut.Lock()
	client.reauthmut.ongoing = nil
	client.reauthmut.Unlock()

	// Report result to all other interested goroutines.
	//
	// This happens in a separate goroutine because another goroutine might have
	// acquired a copy of `client.reauthmut.ongoing` before we cleared it, but not
	// have come around to sending its request. By answering in a goroutine, we
	// can have that goroutine linger until all responseChannels have been sent.
	// When GC has collected all sendings ends of the channel, our receiving end
	// will be closed and the goroutine will end.
	go func() {
		for responseChannel := range messages {
			responseChannel <- err
		}
	}()
	return err
}

// SetDebug for request and response
func (client *ProviderClient) SetDebug(debug bool) {
	client.debug = debug
	log.SetLevel(log.DebugLevel)
}

// EnableGetRetriesOn5XX enables GET retries on 5XX errors with the specified number of attempts and a base interval
// (in seconds) following an exponential backoff with jitter strategy.
func (client *ProviderClient) EnableGetRetriesOn5XX(attempts int, interval int) {
	client.retryGetOn5XX = true
	client.retryGetOn5XXAttempts = attempts
	client.retryGetOn5XXBaseInterval = interval
}

func (client *ProviderClient) IsDebug() bool {
	return client.debug
}

func (client *ProviderClient) debugRequest(request *http.Request) {
	if client.debug {
		dump, err := httputil.DumpRequestOut(request, true)
		if err != nil {
			log.Error(err)
		} else {
			log.SetOutput(os.Stderr)
			log.Debug(string(dump))
			log.SetOutput(os.Stdout)
		}
	}
}

func (client *ProviderClient) debugResponse(response *http.Response) {
	if client.debug {
		dump, err := httputil.DumpResponse(response, true)
		if err != nil {
			log.Error(err)
		} else {
			log.SetOutput(os.Stderr)
			log.Debug(string(dump))
			log.SetOutput(os.Stdout)
		}
	}
}

// RequestOpts customizes the behavior of the provider.Request() method.
type RequestOpts struct {
	// JSONBody, if provided, will be encoded as JSON and used as the body of the HTTP request. The
	// content type of the request will default to "application/json" unless overridden by MoreHeaders.
	// It's an error to specify both a JSONBody and a RawBody.
	JSONBody interface{}
	// RawBody contains an io.Reader that will be consumed by the request directly. No content-type
	// will be set unless one is provided explicitly by MoreHeaders.
	RawBody io.Reader
	// JSONResponse, if provided, will be populated with the contents of the response body parsed as
	// JSON.
	JSONResponse interface{}
	// OkCodes contains a list of numeric HTTP status codes that should be interpreted as success. If
	// the response has a different code, an error will be returned.
	OkCodes []int
	// MoreHeaders specifies additional HTTP headers to be provide on the request. If a header is
	// provided with a blank value (""), that header will be *omitted* instead: use this to suppress
	// the default Accept header or an inferred Content-Type, for example.
	MoreHeaders map[string]string
	// ErrorContext specifies the resource error type to return if an error is encountered.
	// This lets resources override default error messages based on the response status code.
	ErrorContext error
	// ConflictRetryAmount specifies number of retries to perform in case of encountering conflict
	// during performed request
	ConflictRetryAmount int
	// ConflictRetryInterval specifies time (in seconds) between next retry requests
	ConflictRetryInterval int
}

// requestState contains temporary state for a single ProviderClient.Request() call.
type requestState struct {
	// This flag indicates if we have reauthenticated during this request because of a 401 response.
	// It ensures that we don't reauthenticate multiple times for a single request. If we
	// reauthenticate, but keep getting 401 responses with the fresh token, reauthenticating some more
	// will just get us into an infinite loop.
	hasReauthenticated bool
	// This flag indicates if issued request has already been retried. It ensures that we don't
	// end inside an infinite loop.
	hasRetried bool
}

var applicationJSON = "application/json"

// Request performs an HTTP request using the ProviderClient's current HTTPClient. An authentication
// header will automatically be provided.
func (client *ProviderClient) Request(method, url string, options *RequestOpts) (*http.Response, error) {
	return client.doRequest(method, url, options, &requestState{
		hasReauthenticated: false,
		hasRetried:         false,
	})
}

func (client *ProviderClient) doRequest(method, url string, options *RequestOpts, state *requestState) (*http.Response, error) { // nolint: gocyclo
	var body io.Reader
	var contentType *string

	// Derive the content body by either encoding an arbitrary object as JSON, or by taking a provided
	// io.ReadSeeker as-is. Default the content-type to application/json.
	if options.JSONBody != nil {
		if options.RawBody != nil {
			return nil, errors.New("please provide only one of JSONBody or RawBody to gcorecloud.Request()")
		}

		rendered, err := json.Marshal(options.JSONBody)
		if err != nil {
			return nil, err
		}

		body = bytes.NewReader(rendered)
		contentType = &applicationJSON
	}

	if options.RawBody != nil {
		body = options.RawBody
	}

	// Construct the http.Request.
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if client.Context != nil {
		req = req.WithContext(client.Context)
	}

	// Populate the request headers. Apply options.MoreHeaders last, to give the caller the chance to
	// modify or omit any header.
	if contentType != nil {
		req.Header.Set("Content-Type", *contentType)
	}
	req.Header.Set("Accept", applicationJSON)

	// Set the User-Agent header
	req.Header.Set("User-Agent", client.UserAgent.Join())

	if options.MoreHeaders != nil {
		for k, v := range options.MoreHeaders {
			if v != "" {
				req.Header.Set(k, v)
			} else {
				req.Header.Del(k)
			}
		}
	}

	// get latest token from client
	for k, v := range client.AuthenticatedHeaders() {
		req.Header.Set(k, v)
	}

	// Set connection parameter to close the connection immediately when we've got the response
	req.Close = true

	preReqToken := client.AccessToken()

	client.debugRequest(req)

	// Issue the request.
	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	client.debugResponse(resp)

	// Allow default OkCodes if none explicitly set
	okc := options.OkCodes
	if okc == nil {
		okc = defaultOkCodes(method)
	}

	// Validate the HTTP response status.
	var ok bool
	for _, code := range okc {
		if resp.StatusCode == code {
			ok = true
			break
		}
	}

	if !ok {
		body, _ := io.ReadAll(resp.Body)
		err := resp.Body.Close()
		if err != nil {
			log.Error(err)
		}
		respErr := ErrUnexpectedResponseCode{
			URL:      url,
			Method:   method,
			Expected: options.OkCodes,
			Actual:   resp.StatusCode,
			Body:     body,
		}

		errType := options.ErrorContext

		// Handling for all GET requests returning 5xx status codes
		if method == http.MethodGet && resp.StatusCode >= 500 && resp.StatusCode < 600 {
			// If retries on GET requests are enabled, retry requests according to client settings (attempts and interval)
			if client.retryGetOn5XX && client.retryGetOn5XXAttempts > 0 && !state.hasRetried {
				state.hasRetried = true
				for attempt := 1; attempt <= client.retryGetOn5XXAttempts; attempt++ {
					resp, err = client.doRequest(method, url, options, state)
					if err != nil {
						log.Warningf("Retried request failed.\nDetails: %v", err)
					}
					// Validate the HTTP response status.
					for _, code := range okc {
						if resp.StatusCode == code && err == nil {
							return resp, err
						}
					}

					// Exponential backoff with jitter
					// E.g for 3 attempts and interval 2s: 1s, 3s, 7s
					sleepDuration := time.Duration(float64(client.retryGetOn5XXBaseInterval) *
						math.Pow(2, float64(attempt-1)) * (0.5 + rand.Float64()/2))
					time.Sleep(sleepDuration * time.Second)
				}
			}
		}

		switch resp.StatusCode {
		case http.StatusBadRequest:
			err = ErrDefault400{respErr}
			if error400er, ok := errType.(Err400er); ok {
				err = error400er.Error400(respErr)
			}
		case http.StatusUnauthorized:
			if client.ReauthFunc != nil && !state.hasReauthenticated {
				err = client.Reauthenticate(preReqToken)
				if err != nil {
					e := &ErrUnableToReauthenticate{}
					e.ErrOriginal = respErr
					return nil, e
				}
				if options.RawBody != nil {
					if seeker, ok := options.RawBody.(io.Seeker); ok {
						_, err = seeker.Seek(0, 0)
						if err != nil {
							log.Error(err)
						}
					}
				}
				state.hasReauthenticated = true
				resp, err = client.doRequest(method, url, options, state)
				if err != nil {
					switch err := err.(type) {
					case *ErrUnexpectedResponseCode:
						e := &ErrErrorAfterReauthentication{}
						e.ErrOriginal = err
						return nil, e
					default:
						e := &ErrErrorAfterReauthentication{}
						e.ErrOriginal = err
						return nil, e
					}
				}
				return resp, nil
			}
			err = ErrDefault401{respErr}
			if error401er, ok := errType.(Err401er); ok {
				err = error401er.Error401(respErr)
			}
		case http.StatusForbidden:
			err = ErrDefault403{respErr}
			if error403er, ok := errType.(Err403er); ok {
				err = error403er.Error403(respErr)
			}
		case http.StatusNotFound:
			err = ErrDefault404{respErr}
			if error404er, ok := errType.(Err404er); ok {
				err = error404er.Error404(respErr)
			}
		case http.StatusMethodNotAllowed:
			err = ErrDefault405{respErr}
			if error405er, ok := errType.(Err405er); ok {
				err = error405er.Error405(respErr)
			}
		case http.StatusRequestTimeout:
			err = ErrDefault408{respErr}
			if error408er, ok := errType.(Err408er); ok {
				err = error408er.Error408(respErr)
			}
		case http.StatusConflict:
			if options.ConflictRetryAmount > 0 && !state.hasRetried {
				state.hasRetried = true
				for attempt := 1; attempt <= options.ConflictRetryAmount; attempt++ {
					resp, err := client.doRequest(method, url, options, state)

					if err != nil {
						log.Warningf("Retried request fails, trying again in %d seconds.\nDetails: %v", options.ConflictRetryInterval, err)
					}

					// Validate the HTTP response status.
					for _, code := range okc {
						if resp.StatusCode == code && err == nil {
							return resp, err
						}
					}

					time.Sleep(time.Duration(options.ConflictRetryInterval) * time.Second)
				}
			}
			err = ErrDefault409{respErr}
			if error409er, ok := errType.(Err409er); ok {
				err = error409er.Error409(respErr)
			}
		case http.StatusTooManyRequests:
			err = ErrDefault429{respErr}
			if error429er, ok := errType.(Err429er); ok {
				err = error429er.Error429(respErr)
			}
		case http.StatusInternalServerError:
			err = ErrDefault500{respErr}
			if error500er, ok := errType.(Err500er); ok {
				err = error500er.Error500(respErr)
			}
		case http.StatusBadGateway:
			err = ErrDefault502{respErr}
			if error502er, ok := errType.(Err502er); ok {
				err = error502er.Error502(respErr)
			}
		case http.StatusServiceUnavailable:
			err = ErrDefault503{respErr}
			if error503er, ok := errType.(Err503er); ok {
				err = error503er.Error503(respErr)
			}
		case http.StatusGatewayTimeout:
			err = ErrDefault504{respErr}
			if error504er, ok := errType.(Err504er); ok {
				err = error504er.Error504(respErr)
			}
		}

		if err == nil {
			err = respErr
		}

		return resp, err
	}

	// Parse the response body as JSON, if requested to do so.
	if options.JSONResponse != nil {
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				log.Error(err)
			}
		}()
		if err := json.NewDecoder(resp.Body).Decode(options.JSONResponse); err != nil {
			return nil, err
		}
	}

	return resp, nil
}

// ToTokenOptions - TokenOptions from ProviderClient
func (client ProviderClient) ToTokenOptions() TokenOptions {
	return TokenOptions{
		RefreshToken: client.RefreshToken(),
		AccessToken:  client.AccessToken(),
		AllowReauth:  false,
		APIURL:       client.IdentityEndpoint,
	}
}

func defaultOkCodes(method string) []int {
	switch {
	case method == "GET":
		return []int{200}
	case method == "POST":
		return []int{201, 202, 200}
	case method == "PUT":
		return []int{201, 202}
	case method == "PATCH":
		return []int{200, 202, 204}
	case method == "DELETE":
		return []int{200, 202, 204}
	}

	return []int{}
}

// NewProviderClient - Default constructor
func NewProviderClient() *ProviderClient {
	client := new(ProviderClient)
	return client
}
