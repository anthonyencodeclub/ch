package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/alecthomas/kong"

	"github.com/anthonyencodeclub/ch/internal/errfmt"
	"github.com/anthonyencodeclub/ch/internal/outfmt"
	"github.com/anthonyencodeclub/ch/internal/ui"
)

const (
	colorAuto  = "auto"
	colorNever = "never"
)

// RootFlags are flags available on every command.
type RootFlags struct {
	Color   string `help:"Color output: auto|always|never" default:"auto"`
	JSON    bool   `help:"Output JSON to stdout (best for scripting)" default:"false"`
	Plain   bool   `help:"Output stable, parseable text to stdout (no colors)" default:"false"`
	Verbose bool   `help:"Enable verbose logging"`
}

// CLI is the top-level command tree.
type CLI struct {
	RootFlags `embed:""`

	Version    kong.VersionFlag `help:"Print version and exit"`

	Auth       AuthCmd       `cmd:"" help:"Manage API key authentication"`
	Setup      SetupCmd      `cmd:"" help:"Set up a new company (interactive guided flow)"`
	Company    CompanyCmd    `cmd:"" help:"Company profile and registered office"`
	Search     SearchCmd     `cmd:"" help:"Search companies, officers, and disqualified officers"`
	Officers   OfficersCmd   `cmd:"" help:"List and view company officers"`
	Filing     FilingCmd     `cmd:"" help:"Filing history"`
	PSC        PSCCmd        `cmd:"" help:"Persons with significant control"`
	Charges    ChargesCmd    `cmd:"" help:"Company charges (mortgages/securities)"`
	Insolvency InsolvencyCmd `cmd:"" help:"Insolvency information"`
	File       FileCmd       `cmd:"" help:"File changes (registered address, email) — requires OAuth2 login"`
	VersionCmd VersionCmd    `cmd:"" name:"version" help:"Print version"`
}

type exitPanic struct{ code int }

// Execute parses and runs the CLI.
func Execute(args []string) (err error) {
	parser, cli, err := newParser()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			if ep, ok := r.(exitPanic); ok {
				if ep.code == 0 {
					err = nil
					return
				}
				err = &ExitError{Code: ep.code, Err: errors.New("exited")}
				return
			}
			panic(r)
		}
	}()

	kctx, err := parser.Parse(args)
	if err != nil {
		parsedErr := wrapParseError(err)
		_, _ = fmt.Fprintln(os.Stderr, errfmt.Format(parsedErr))
		return parsedErr
	}

	logLevel := slog.LevelWarn
	if cli.Verbose {
		logLevel = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	})))

	mode, err := outfmt.FromFlags(cli.JSON, cli.Plain)
	if err != nil {
		return &ExitError{Code: 2, Err: err}
	}

	ctx := context.Background()
	ctx = outfmt.WithMode(ctx, mode)

	uiColor := cli.Color
	if outfmt.IsJSON(ctx) || outfmt.IsPlain(ctx) {
		uiColor = colorNever
	}

	u, err := ui.New(ui.Options{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Color:  uiColor,
	})
	if err != nil {
		return err
	}
	ctx = ui.WithUI(ctx, u)

	kctx.BindTo(ctx, (*context.Context)(nil))
	kctx.Bind(&cli.RootFlags)

	err = kctx.Run()
	if err == nil {
		return nil
	}

	if u := ui.FromContext(ctx); u != nil {
		u.Error(errfmt.Format(err))
		return err
	}
	_, _ = fmt.Fprintln(os.Stderr, errfmt.Format(err))
	return err
}

func wrapParseError(err error) error {
	if err == nil {
		return nil
	}
	var parseErr *kong.ParseError
	if errors.As(err, &parseErr) {
		return &ExitError{Code: 2, Err: parseErr}
	}
	return err
}

func newParser() (*kong.Kong, *CLI, error) {
	envMode := outfmt.FromEnv()
	_ = envMode // used for defaults if needed

	cli := &CLI{}
	parser, err := kong.New(
		cli,
		kong.Name("ch"),
		kong.Description("Companies House CLI — search, inspect, and explore UK company data from your terminal."),
		kong.Vars{"version": VersionString()},
		kong.Writers(os.Stdout, os.Stderr),
		kong.Exit(func(code int) { panic(exitPanic{code: code}) }),
	)
	if err != nil {
		return nil, nil, err
	}
	return parser, cli, nil
}
