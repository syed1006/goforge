package generators

// Pinned module versions used by every generator. Bumping these is the only
// place a scaffold's dependency tree should change — keep them explicit so
// every scaffold of the same goforge build produces the same go.sum.
const (
	VersionChi          = "v5.1.0"
	VersionGin          = "v1.10.0"
	VersionFiber        = "v2.52.5"
	VersionEcho         = "v4.12.0"
	VersionFiberAdaptor = "v2.2.1"

	VersionPgx       = "v5.7.1"
	VersionMySQL     = "v1.8.1"
	VersionSQLite    = "v1.33.1"
	VersionMongo     = "v1.17.1"
	VersionRedis     = "v9.7.0"

	VersionGRPC   = "v1.68.0"
	VersionGqlgen = "v0.17.55"
)
