// Scry Info.  All rights reserved.
// license that can be found in the license file.

package dot

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"regexp"
	"strings"

	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/approle"

	"github.com/scryinfo/scryg/sutils/skit"
)

const (

	//SconfigTypeID scofig dot type id
	SconfigTypeID = "484ef01d-3c04-4517-a643-2d776a9ae758"
)

var reVar = regexp.MustCompile(`^\${(\w+)}$`)

var envMap map[string]interface{}

type StringFromEnv string

func (e *StringFromEnv) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return err
	}
	if match := reVar.FindStringSubmatch(s); len(match) > 0 {
		if envMap != nil {
			if value, ok := envMap[match[1]].(string); ok {
				*e = StringFromEnv(value)
			}
		} else {
			*e = StringFromEnv(os.Getenv(match[1]))
		}
	} else {
		*e = StringFromEnv(s)
	}
	return nil
}

//SConfig config belongs to one component Dot, but it is so basic, every Dot need it, so define it in dot.go file
//S represents scryinfo config this name is used frequently, so add s to distinguish it
type SConfig interface {
	//RootPath root path
	RootPath()
	//Config file path
	ConfigPath() string
	//Without path, only file name
	ConfigFile() string
	//Whether key existing
	Key(key string) bool
	//If no config or config is empty, return nil
	Map() (m map[string]interface{}, err error)

	//Priorly bring data from RAM and operate, if RAM is nil then check whether original config file existing
	Unmarshal(s interface{}) error
	//Analyze key is corresponding type
	UnmarshalKey(key string, obj interface{}) error

	Marshal(data []byte) error

	//If no corresponding key or data type cannot be converted, must pay attention to default value,
	//so add "Def" before function to notify default value
	DefInterface(key string, def interface{}) interface{}
	DefArray(key string, def []interface{}) []interface{}
	DefMap(key string, def map[string]interface{}) map[string]interface{}
	DefString(key string, def string) string
	DefInt32(key string, def int32) int32
	DefUint32(key string, def uint32) uint32
	DefInt64(key string, def int64) int64
	DefUint64(key string, def uint64) uint64
	DefBool(key string, def bool) bool
	DefFloat32(key string, def float32) float32
	DefFloat64(key string, def float64) float64
}

func readIDFromFile(file string) (string, error) {
	secretIDFile, err := os.Open(file)
	if err != nil {
		return "", fmt.Errorf("unable to open file containing secret ID: %w", err)
	}
	defer secretIDFile.Close()

	limitedReader := io.LimitReader(secretIDFile, 1000)
	secretIDBytes, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", fmt.Errorf("unable to read secret ID: %w", err)
	}

	secretIDValue := strings.TrimSuffix(string(secretIDBytes), "\n")

	return secretIDValue, nil
}

// Fetches a key-value secret (kv-v2) after authenticating via AppRole.
func GetSecretWithAppRole(keypath string, vaultAdd string) (map[string]interface{}, error) {
	config := vault.DefaultConfig() // modify for more granular configuration

	config.Address = vaultAdd

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault client: %w", err)
	}
	// A combination of a Role ID and Secret ID is required to log in to Vault
	// with an AppRole.
	// First, let's get the role ID given to us by our Vault administrator.
	//roleID := os.Getenv("APPROLE_ROLE_ID")
	//roleIdFile := os.Getenv("VAULT_ROLE_ID_FILE")
	roleID, err := readIDFromFile(".vault_role_id")
	if err != nil {
		roleID, err = readIDFromFile("/run/secrets/vault_role_id")
		if err != nil {
			return nil, fmt.Errorf("no role ID was provided in APPROLE_ROLE_ID env var")
		}
	}

	// The Secret ID is a value that needs to be protected, so instead of the
	// app having knowledge of the secret ID directly, we have a trusted orchestrator (https://learn.hashicorp.com/tutorials/vault/secure-introduction?in=vault/app-integration#trusted-orchestrator)
	// give the app access to a short-lived response-wrapping token (https://www.vaultproject.io/docs/concepts/response-wrapping).
	// Read more at: https://learn.hashicorp.com/tutorials/vault/approle-best-practices?in=vault/auth-methods#secretid-delivery-best-practices
	//secretIdFile := os.Getenv("VAULT_SECRET_ID_FILE")
	secretIdFile := ".vault_secret_id"
	_, err = readIDFromFile(secretIdFile)
	if err != nil {
		secretIdFile = "/run/secrets/vault_secret_id"
		_, err = readIDFromFile(secretIdFile)
		if err != nil {
			return nil, fmt.Errorf("no role ID was provided in APPROLE_SECRET_ID ")
		}
	}
	secretID := &auth.SecretID{FromFile: secretIdFile}

	appRoleAuth, err := auth.NewAppRoleAuth(
		roleID,
		secretID,
		//auth.WithWrappingToken(), // Only required if the secret ID is response-wrapped.
	)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize AppRole auth method: %w", err)
	}

	authInfo, err := client.Auth().Login(context.Background(), appRoleAuth)
	if err != nil {
		return nil, fmt.Errorf("unable to login to AppRole auth method: %w", err)
	}
	if authInfo == nil {
		return nil, fmt.Errorf("no auth info was returned after login")
	}

	// get secret from the default mount path for KV v2 in dev mode, "secret"
	secret, err := client.KVv2("secret").Get(context.Background(), keypath)
	if err != nil {
		return nil, fmt.Errorf("unable to read secret: %w", err)
	}

	return secret.Data, nil
}

//UnMarshalConfig unmarshal config
func UnMarshalConfig(conf []byte, obj interface{}) (err error) {
	err = nil
	if conf != nil {
		err = yaml.Unmarshal(conf, obj)
		//err = json.Unmarshal(conf, obj)
	} else {
		err = SError.Parameter
	}
	return err
}

//MarshalConfig marshal config
func MarshalConfig(lconf *LiveConfig) (conf []byte, err error) {
	conf = nil
	err = nil

	if lconf != nil {
		if !skit.IsNil(lconf.Config) {
			conf, ok := lconf.Config.(map[string]interface{})
			if ok {
				keypath, ok := conf["keypath"].(string)
				if ok {
					if keypath != "" {
						VAULT_ADDR := os.Getenv("VAULT_ADDR")
						if VAULT_ADDR != "" {
							//Logger().Infoln("marshal ", zap.String("keypath", keypath))
							keys, err := GetSecretWithAppRole(keypath, VAULT_ADDR)
							if err == nil {
								envMap = keys
								Logger().Infoln("envMap ", zap.Any("envMap", envMap))
								//for key, value := range keys {
								//	if value == "true" || value == "false" {
								//		switch value {
								//		case "true":
								//			conf[key] = true
								//		case "false":
								//			conf[key] = false
								//		}
								//	} else {
								//		conf[key] = value
								//	}
								//}
							} else {
								Logger().Errorln("Get keypath failed", zap.String("keypath", keypath))
							}
						}
						Logger().Infoln("marshal ok", zap.String("keypath", keypath))
					}
				}
			}

			return yaml.Marshal(lconf.Config)
			//return json.Marshal(lconf.Config)
		}
	}
	return conf, err
}
