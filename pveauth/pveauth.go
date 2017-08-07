package pveauth

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/andrexus/goproxmox/pveauth/internal"
)

type Config struct {
	// Username
	Username string

	// Password
	Password string

	TicketURL string
}

// A TicketSource is anything that can return a ticket.
type TicketSource interface {
	// Ticket returns a ticket or an error.
	// Ticket must be safe for concurrent use by multiple goroutines.
	// The returned Ticket must not be modified.
	Ticket() (*Ticket, error)
}

// PasswordCredentialsTicket converts a resource owner username and password
// pair into a ticket.
//
// The HTTP client to use is derived from the context.
// If nil, http.DefaultClient is used.
func (c *Config) PasswordCredentialsTicket(ctx context.Context) (*Ticket, error) {
	return retrieveTicket(ctx, c.Username, c.Password, c.TicketURL)
}

// Client returns an HTTP client using the provided ticket.
// The ticket will auto-refresh as necessary. The underlying
// HTTP transport will be obtained using the provided context.
// The returned client and its Transport should not be modified.
func (c *Config) Client(ctx context.Context, t *Ticket) *http.Client {
	return NewClient(ctx, c.TicketSource(ctx, t))
}

// TicketSource returns a TicketSource that returns t until t expires,
// automatically refreshing it as necessary using the provided context.
//
// Most users will use Config.Client instead.
func (c *Config) TicketSource(ctx context.Context, t *Ticket) TicketSource {
	tkr := &ticketRefresher{
		ctx:  ctx,
		conf: c,
	}
	if t != nil {
		tkr.refreshTicket = t.Ticket
	}
	return &reuseTicketSource{
		t:   t,
		new: tkr,
	}
}

// ticketRefresher is a TicketSource that makes "grant_type"=="refresh_ticket"
// HTTP requests to renew a ticket using a RefreshTicket.
type ticketRefresher struct {
	ctx           context.Context // used to get HTTP requests
	conf          *Config
	refreshTicket string
}

// WARNING: Ticket is not safe for concurrent access, as it
// updates the ticketRefresher's refreshTicket field.
// Within this package, it is used by reuseTicketSource which
// synchronizes calls to this method with its own mutex.
func (tf *ticketRefresher) Ticket() (*Ticket, error) {
	if tf.refreshTicket == "" {
		return nil, errors.New("pveauth: ticket expired and refresh ticket is not set")
	}

	tk, err := retrieveTicket(tf.ctx, tf.conf.Username, tf.conf.Password, tf.conf.TicketURL)

	if err != nil {
		return nil, err
	}
	if tf.refreshTicket != tk.Ticket {
		tf.refreshTicket = tk.Ticket
	}
	return tk, err
}

// reuseTicketSource is a TicketSource that holds a single ticket in memory
// and validates its expiry before each call to retrieve it with
// Ticket. If it's expired, it will be auto-refreshed using the
// new TicketSource.
type reuseTicketSource struct {
	new TicketSource // called when t is expired.

	mu sync.Mutex // guards t
	t  *Ticket
}

// Ticket returns the current ticket if it's still valid, else will
// refresh the current ticket (using r.Context for HTTP client
// information) and return the new one.
func (s *reuseTicketSource) Ticket() (*Ticket, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.t.Valid() {
		return s.t, nil
	}
	t, err := s.new.Ticket()
	if err != nil {
		return nil, err
	}
	s.t = t
	return t, nil
}

// HTTPClient is the context key to use with golang.org/x/net/context's
// WithValue function to associate an *http.Client value with a context.
var HTTPClient internal.ContextKey

// NewClient creates an *http.Client from a Context and TicketSource.
// The returned client is not valid beyond the lifetime of the context.
//
// As a special case, if src is nil, a non-OAuth2 client is returned
// using the provided context. This exists to support related OAuth2
// packages.
func NewClient(ctx context.Context, src TicketSource) *http.Client {
	if src == nil {
		c, err := internal.ContextClient(ctx)
		if err != nil {
			return &http.Client{Transport: internal.ErrorTransport{Err: err}}
		}
		return c
	}
	return &http.Client{
		Transport: &Transport{
			Base:   internal.ContextTransport(ctx),
			Source: ReuseTicketSource(nil, src),
		},
	}
}

// ReuseTicketSource returns a TicketSource which repeatedly returns the
// same ticket as long as it's valid, starting with t.
// When its cached ticket is invalid, a new ticket is obtained from src.
//
// ReuseTicketSource is typically used to reuse tickets from a cache
// (such as a file on disk) between runs of a program, rather than
// obtaining new tickets unnecessarily.
//
// The initial ticket t may be nil, in which case the TicketSource is
// wrapped in a caching version if it isn't one already. This also
// means it's always safe to wrap ReuseTicketSource around any other
// TicketSource without adverse effects.
func ReuseTicketSource(t *Ticket, src TicketSource) TicketSource {
	// Don't wrap a reuseTicketSource in itself. That would work,
	// but cause an unnecessary number of mutex operations.
	// Just build the equivalent one.
	if rt, ok := src.(*reuseTicketSource); ok {
		if t == nil {
			// Just use it directly.
			return rt
		}
		src = rt.new
	}
	return &reuseTicketSource{
		t:   t,
		new: src,
	}
}
