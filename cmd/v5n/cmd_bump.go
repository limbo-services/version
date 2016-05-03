package main

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/cloud/datastore"
	"gopkg.in/alecthomas/kingpin.v2"
	"limbo.services/version"
)

func bump(app *kingpin.Application, appName, segmentName, gcpProjectID string) {
	ctx := context.Background()
	ctx = datastore.WithNamespace(ctx, "V5nStore")
	ds, err := datastore.NewClient(ctx, gcpProjectID)
	app.FatalIfError(err, "v5n")

	var (
		appKey        = datastore.NewKey(ctx, "Application", appName, 0, nil)
		reportVersion *version.Info
	)

	_, err = ds.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var info remoteVersionInfo

		err := tx.Get(appKey, &info)
		if err == datastore.ErrNoSuchEntity {
			info.Name = appName
			info.Version = "0.0.0"
			err = nil
		}
		if err != nil {
			return err
		}

		vsn, err := version.ParseSemver(info.Version)
		if err != nil {
			return err
		}

		switch segmentName {
		case "major":
			vsn = version.Bump(vsn, version.Major)
			reportVersion = vsn
		case "minor":
			vsn = version.Bump(vsn, version.Minor)
			reportVersion = vsn
		case "patch":
			vsn = version.Bump(vsn, version.Patch)
			reportVersion = vsn
		case "rc":
			vsn = version.Bump(vsn, version.ReleaseCandidate)
			reportVersion = vsn
		case "final+major":
			reportVersion = version.Bump(vsn, version.Final)
			vsn = version.Bump(vsn, version.Major)
		case "final+minor":
			reportVersion = version.Bump(vsn, version.Final)
			vsn = version.Bump(vsn, version.Minor)
		case "final+patch":
			reportVersion = version.Bump(vsn, version.Final)
			vsn = version.Bump(vsn, version.Patch)
		default:
			return errors.New("invalid bump type")
		}

		info.Version = vsn.Semver()
		_, err = tx.Put(appKey, &info)
		return err
	})
	app.FatalIfError(err, "v5n")

	fmt.Fprintf(os.Stdout, "%s\n", reportVersion.Semver())
}

type remoteVersionInfo struct {
	Name    string
	Version string
}
