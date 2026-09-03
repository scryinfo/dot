package daokeys

var (
	PrefixAuthCode       = []byte("oidc:auth_code:")
	PrefixAuthRequest    = []byte("oidc:auth_request:")
	PrefixUser           = []byte("oidc:user:")
	PrefixIdentity       = []byte("oidc:identity:")
	PrefixUserIdentities = []byte("oidc:user_identities:")
	PrefixToken          = []byte("oidc:token:")
	PrefixRefreshToken   = []byte("oidc:refresh_token:")
	PrefixOidcClient     = []byte("oidc:client:")
	PrefixClientStatus   = []byte("oidc:client_status:")

	// PrefixMessage     = []byte("admin:message:")
	// PrefixMessageBody = []byte("admin:message_body:")

	SettingOg = []byte("admin:setting:og:")
	SettingTg = []byte("admin:setting:tg:")
)
