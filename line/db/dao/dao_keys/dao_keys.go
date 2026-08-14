package daokeys

var (
	PrefixBenchData = []byte("bench:data:")
	// the maximum player id
	KeyPlayersMeta = []byte("admin:player_meta:")
	PrefixPlayer   = []byte("admin:player:")

	PrefixBanPlayers = []byte("admin:ban_players:")

	PrefixAssets     = []byte("admin:assets:")
	PrefixAssetsBody = []byte("admin:assets_body:")

	PrefixMessage     = []byte("admin:message:")
	PrefixMessageBody = []byte("admin:message_body:")

	PrefixAlltypeMeta = []byte("admin:message_all:")
	PrefixMessageAuto = []byte("admin:message_auto:")

	PrefixWithrowDataBody = []byte("admin:withdraw_body:")
	PrefixWithrowData     = []byte("admin:withdraw:")

	PrefixUser = []byte("admin:user:")

	PrefixInvitationCode = []byte("admin:invitation:")
	PrefixOgPlayer       = []byte("admin:ogplayer:")
	PrefixOgQuery        = []byte("admin:ogquery:")

	SettingOg = []byte("admin:setting:og:")
	SettingTg = []byte("admin:setting:tg:")
)
