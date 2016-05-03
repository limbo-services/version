package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
	"limbo.services/version"
)

func main() {
	var (
		binPath        string
		versionArg     string
		commitArg      string
		authorArg      string
		releaseDateArg string
		appName        string
		segmentName    string
		gcpProjectID   string
	)

	app := kingpin.New("v5n", "Go program version manager").
		Version(version.Get().String()).
		Author(version.Get().ReleasedBy)

	applyCmd := app.Command("apply", "apply version info to a go binary")
	applyCmd.Arg("binary", "binary to patch").Required().ExistingFileVar(&binPath)
	applyCmd.Arg("version", "version to apply").Required().StringVar(&versionArg)
	applyCmd.Flag("commit", "commit hash").StringVar(&commitArg)
	applyCmd.Flag("author", "author info").StringVar(&authorArg)
	applyCmd.Flag("date", "release date").StringVar(&releaseDateArg)

	store := app.Command("store", "manage the version store")
	store.Flag("gcp-project", "GCP project id").
		PlaceHolder("GCP_PROJECT_ID").Envar("GCP_PROJECT_ID").
		Required().StringVar(&gcpProjectID)

	bumpCmd := store.Command("bump", "bump a version in the store")
	bumpCmd.Arg("app", "name of the application").Required().StringVar(&appName)
	bumpCmd.Arg("segment", "segment to bump").Required().EnumVar(&segmentName,
		"major", "minor", "patch", "rc", "final+major", "final+minor", "final+patch")

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case applyCmd.FullCommand():
		apply(app, binPath, versionArg, authorArg, commitArg, releaseDateArg)
	case bumpCmd.FullCommand():
		bump(app, appName, segmentName, gcpProjectID)
	}
}
