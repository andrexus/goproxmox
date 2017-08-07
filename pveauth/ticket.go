package pveauth

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/andrexus/goproxmox/pveauth/internal"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

// expiryDelta determines how earlier a ticket should be considered
// expired than its actual expiration time. It is used to avoid late
// expirations due to client-server time mismatches.
const expiryDelta = 10 * time.Second
const ticketValidityTime = 2 * time.Hour

type ticketRoot struct {
	Ticket Ticket `json:"data"`
}

// Ticket represents the crendentials used to authorize
// the requests to access protected resources on the PVE server.
//
// Most users of this package should not access fields of Ticket
// directly.
type Ticket struct {
	// Ticket is the ticket that authorizes and authenticates
	// the requests.
	Ticket string `json:"ticket"`

	// CSRFPreventionToken is used for all write requests.
	CSRFPreventionToken string `json:"CSRFPreventionToken"`

	// Username
	Username string `json:"username"`
}

// SetAuthHeader sets the Authorization header to r using the access
// ticket in t.
//
// This method is unnecessary when using Transport or an HTTP Client
// returned by this package.
func (t *Ticket) SetAuthHeader(r *http.Request) {
	r.Header.Set("Cookie", fmt.Sprintf("PVEAuthCookie=%s", t.Ticket))
	if r.Method == "POST" || r.Method == "PUT" {
		r.Header.Add("CSRFPreventionToken", t.CSRFPreventionToken)
	}
}

// expired reports whether the ticket is expired.
// t must be non-nil.
func (t *Ticket) expired() bool {
	return time.Now().Add(ticketValidityTime).Add(-expiryDelta).Before(time.Now())
}

// Valid reports whether t is non-nil, has an AccessTicket, and is not expired.
func (t *Ticket) Valid() bool {
	return t != nil && t.Ticket != "" && !t.expired()
}

func retrieveTicket(ctx context.Context, username, password, ticketURL string) (*Ticket, error) {
	hc, err := internal.ContextClient(ctx)
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	v.Add("username", username)
	v.Add("password", password)

	req, err := http.NewRequest("POST", ticketURL, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	r, err := ctxhttp.Do(ctx, hc, req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("pveauth: cannot fetch ticket: %v", err)
	}
	if code := r.StatusCode; code < 200 || code > 299 {
		return nil, fmt.Errorf("pveauth: cannot fetch ticket: %v\nResponse: %s", r.Status, body)
	}

	var ticket *Ticket

	var tr ticketRoot
	if err = json.Unmarshal(body, &tr); err != nil {
		return nil, err
	}
	ticket = &Ticket{
		Ticket:              tr.Ticket.Ticket,
		CSRFPreventionToken: tr.Ticket.CSRFPreventionToken,
		Username:            tr.Ticket.Username,
	}

	return ticket, nil
}
