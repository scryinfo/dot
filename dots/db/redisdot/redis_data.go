package redisdot

const (
	VersionControlChannelName = "version_control"
	KeyWithVersionTemplate = "%s:v%d"

	GetVersionNotExistFlag = iota - 2
	GetVersionNotExistAndRegisterFlag
)
