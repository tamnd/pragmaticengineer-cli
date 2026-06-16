package pragmaticengineer

import (
	"context"

	"github.com/tamnd/any-cli/kit"
	"github.com/tamnd/any-cli/kit/errs"
)

// domain.go exposes pragmaticengineer as a kit Domain: a driver that a multi-domain
// host (ant) enables with a single blank import,
//
//	import _ "github.com/tamnd/pragmaticengineer-cli/pragmaticengineer"
//
// exactly as a database/sql program enables a driver with `import _
// "github.com/lib/pq"`. The init below registers it; the host then dereferences
// pragmaticengineer:// URIs by routing to the operations Register installs.
func init() { kit.Register(Domain{}) }

// Domain is the pragmaticengineer driver. It carries no state; the per-run client is
// built by the factory Register hands kit.
type Domain struct{}

// Info describes the scheme, the hostnames a pasted link is matched against, and
// the identity reused for the binary's help and version.
func (Domain) Info() kit.DomainInfo {
	return kit.DomainInfo{
		Scheme: "pragmaticengineer",
		Hosts:  []string{Host},
		Identity: kit.Identity{
			Binary: "pragmaticengineer",
			Short:  "Browse The Pragmatic Engineer newsletter from the command line.",
			Long: `Browse The Pragmatic Engineer newsletter from the command line.

pragmaticengineer reads public data from The Pragmatic Engineer Substack over plain
HTTPS, shapes it into clean records, and prints output that pipes into the rest
of your tools. No API key, nothing to run alongside it.`,
			Site: Host,
			Repo: "https://github.com/tamnd/pragmaticengineer-cli",
		},
	}
}

// Register installs the client factory and every operation onto app.
func (Domain) Register(app *kit.App) {
	app.SetClient(newClientFactory)

	kit.Handle(app, kit.OpMeta{Name: "top", Group: "read", List: true,
		Summary: "List recent posts from The Pragmatic Engineer",
		Args:    []kit.Arg{}}, topPosts)

	kit.Handle(app, kit.OpMeta{Name: "export", Group: "read", List: true,
		Summary: "Export all posts as JSONL"}, exportPosts)

	kit.Handle(app, kit.OpMeta{Name: "info", Group: "read", Single: true,
		Summary: "Show newsletter statistics"}, getNewsletterInfo)
}

// newClientFactory builds the client from the host-resolved config.
func newClientFactory(_ context.Context, cfg kit.Config) (any, error) {
	dcfg := DefaultConfig()
	if cfg.UserAgent != "" {
		dcfg.UserAgent = cfg.UserAgent
	}
	if cfg.Rate > 0 {
		dcfg.Rate = cfg.Rate
	}
	if cfg.Retries > 0 {
		dcfg.Retries = cfg.Retries
	}
	if cfg.Timeout > 0 {
		dcfg.Timeout = cfg.Timeout
	}
	return NewClient(dcfg), nil
}

// --- inputs ---

type topInput struct {
	Limit  int     `kit:"flag" help:"number of posts to fetch"`
	Offset int     `kit:"flag" help:"offset for pagination"`
	Client *Client `kit:"inject"`
}

type exportInput struct {
	Client *Client `kit:"inject"`
}

type infoInput struct {
	Client *Client `kit:"inject"`
}

// --- handlers ---

func topPosts(ctx context.Context, in topInput, emit func(*Post) error) error {
	limit := in.Limit
	if limit <= 0 {
		limit = 25
	}
	posts, err := in.Client.Top(ctx, limit, in.Offset)
	if err != nil {
		return mapErr(err)
	}
	for _, p := range posts {
		if err := emit(p); err != nil {
			return err
		}
	}
	return nil
}

func exportPosts(ctx context.Context, in exportInput, emit func(*Post) error) error {
	posts, err := in.Client.AllPosts(ctx)
	if err != nil {
		return mapErr(err)
	}
	for _, p := range posts {
		if err := emit(p); err != nil {
			return err
		}
	}
	return nil
}

func getNewsletterInfo(ctx context.Context, in infoInput, emit func(*Info) error) error {
	info, err := in.Client.Stats(ctx)
	if err != nil {
		return mapErr(err)
	}
	return emit(info)
}

// Classify turns any accepted input into (type, id). Only "top" is defined for now.
func (Domain) Classify(input string) (uriType, id string, err error) {
	return "", "", errs.Usage("unrecognized pragmaticengineer reference: %q", input)
}

// Locate is the inverse: the live https URL for a (type, id).
func (Domain) Locate(uriType, id string) (string, error) {
	return "", errs.Usage("pragmaticengineer has no resource type %q", uriType)
}

// mapErr converts a library error into the kit error kind.
func mapErr(err error) error {
	return err
}
