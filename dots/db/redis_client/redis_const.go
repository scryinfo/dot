package redis_client

const (
	VersionControlChannelName = "version_control"

	KeyWithVersionTemplate   = `v%s:%s`
	KeyWithVersionTemplateRE = `v(\w*):(\w*)`

	CurrentVersionPrefix  = "CV:"
	AllVersionsListPrefix = "AVsL:"

	KeySplitChar = ":"
)
