package main

import (
	"crypto/ecdh"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/scryinfo/dot/line/sconfig"
	"github.com/scryinfo/scryg/sutils/sfile"
)

const privateKeyFileExtension = "pri" // dot以常量的形式固定

func main() {
	log.Println("> Dotenc Start:\n" +
		"private key: " + generateFileBaseName + "." + privateKeyFileExtension + "\n" +
		"plaintext  config file: " + plaintextConfigFileName + "\n" +
		"ciphertext config file: " + generateFileBaseName + "." + configFileExtension + "\n")
	defer log.Println("> Dotenc Done.")

	pubKey := getPublicKey()
	encryptConfigFile(pubKey)
}

func getPublicKey() (pubKey *ecdh.PublicKey) {
	privKeyFileName := fmt.Sprintf("%s.%s", generateFileBaseName, privateKeyFileExtension)

	if sfile.ExistFile(privKeyFileName) { // 使用已存在的私钥文件
		privBytes, err := os.ReadFile(privKeyFileName)
		if err != nil {
			log.Fatalln("read private key file failed, error: ", err)
		}

		privKey, err := ecdh.X25519().NewPrivateKey(privBytes)
		if err != nil {
			log.Fatalln("deserialize private key failed, error: ", err)
		}

		pubKey = privKey.PublicKey()
	} else { // 不存在目标私钥文件，创建新的
		privKey, err := ecdh.X25519().GenerateKey(nil)
		if err != nil {
			log.Fatalln("generate private key failed, error: ", err)
		}

		err = os.WriteFile(privKeyFileName, privKey.Bytes(), 0644)
		if err != nil {
			log.Fatalln("write private key file failed, error: ", err)
		}

		pubKey = privKey.PublicKey()
	}

	return
}

func encryptConfigFile(pubKey *ecdh.PublicKey) {
	configBytes, err := os.ReadFile(plaintextConfigFileName)
	if err != nil {
		log.Fatalln("read config file failed, error: ", err)
	}

	ciphertext, err := sconfig.EncriptionFile(configBytes, pubKey)
	if err != nil {
		log.Fatalln("encrypt config file failed, error: ", err)
	}
	ciphertextBase64 := base64.StdEncoding.EncodeToString(ciphertext)

	ciphertextConfigFileName := fmt.Sprintf("%s.%s", generateFileBaseName, configFileExtension)
	err = os.WriteFile(ciphertextConfigFileName, []byte(ciphertextBase64), 0644)
	if err != nil {
		log.Fatalln("write encrypt file failed, error: ", err)
	}
}
