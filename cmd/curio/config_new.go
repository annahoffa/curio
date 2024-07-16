package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/curio/deps"

	"github.com/filecoin-project/lotus/api"
	cliutil "github.com/filecoin-project/lotus/cli/util"
	"github.com/filecoin-project/lotus/node/repo"
)

var configNewCmd = &cli.Command{
	Name:      "new-cluster",
	Usage:     "Create new configuration for a new cluster",
	ArgsUsage: "[SP actor address...]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "repo",
			EnvVars: []string{"LOTUS_PATH"},
			Hidden:  true,
			Value:   "~/.lotus",
		},
	},
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 1 {
			return xerrors.New("must specify at least one SP actor address. Use 'lotus-shed miner create' or use 'curio guided-setup'")
		}

		ctx := cctx.Context
		depnds, err := deps.GetDepsCLI(ctx, cctx)
		if err != nil {
			return xerrors.Errorf("connecting to full node: %w", err)
		}
		db, chain := depnds.DB, depnds.Chain

		ainfo, err := cliutil.GetAPIInfo(cctx, repo.FullNode)
		if err != nil {
			return xerrors.Errorf("could not get API info for FullNode: %w", err)
		}

		token, err := chain.AuthNew(ctx, api.AllPermissions)
		if err != nil {
			return err
		}

		return deps.CreateMinerConfig(ctx, chain, db, cctx.Args().Slice(), fmt.Sprintf("%s:%s", string(token), ainfo.Addr))
	},
}
