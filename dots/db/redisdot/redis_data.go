package redisdot

const (
	VersionControlChannelName = "version_control"
	KeyWithVersionTemplate = "%s:v%d"
	RegisteredKeySuffix = ":currentVersion"
	RegisteredKeysListSuffix = ":versions"

	GetVersionNotExistFlag = iota - 2
	GetVersionNotExistAndRegisterFlag
)
