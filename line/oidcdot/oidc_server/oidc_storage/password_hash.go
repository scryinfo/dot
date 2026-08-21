package oidc_storage

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type _PasswordHash struct{}

var PasswordHash = _PasswordHash{}

// Argon2id 参数配置（OWASP 推荐参数）
type params struct {
	memory      uint32 // 内存消耗 (KiB)，64MB = 64 * 1024
	iterations  uint32 // 迭代次数 (Time)
	parallelism uint8  // 并行度 (Threads)
	saltLength  uint32 // Salt 字节长度
	keyLength   uint32 // 生成的 Hash 字节长度
}

var defaultParams = &params{
	memory:      64 * 1024,
	iterations:  1,
	parallelism: 4,
	saltLength:  16,
	keyLength:   32,
}

// 1. 生成 Hash 字符串（可直接保存到数据库）
func (n _PasswordHash) HashPassword(password string) (string, error) {
	// 1.1 生成随机 Salt
	salt := make([]byte, defaultParams.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// 1.2 使用 argon2.IDKey 计算 Hash
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		defaultParams.iterations,
		defaultParams.memory,
		defaultParams.parallelism,
		defaultParams.keyLength,
	)

	// 1.3 编码为标准的 PHC 格式字符串：$argon2id$v=19$m=65536,t=1,p=4$<salt>$<hash>
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		defaultParams.memory,
		defaultParams.iterations,
		defaultParams.parallelism,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}

// 2. 校验密码与 Hash 是否匹配
func (n _PasswordHash) ComparePasswordAndHash(password, encodedHash string) (bool, error) {
	// 2.1 解析数据库中保存的 PHC 字符串
	p, salt, hash, err := n.decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// 2.2 用解出的参数和 Salt 重新计算输入密码的 Hash
	otherHash := argon2.IDKey(
		[]byte(password),
		salt,
		p.iterations,
		p.memory,
		p.parallelism,
		p.keyLength,
	)

	// 2.3 常数时间比较（防止时序攻击 / Timing Attack）
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

// 辅助函数：解析 PHC 格式字符串
func (n _PasswordHash) decodeHash(encodedHash string) (p *params, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, errors.New("hash 格式不正确")
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil || version != argon2.Version {
		return nil, nil, nil, errors.New("不支持的 argon2 版本")
	}

	p = &params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, errors.New("参数解析失败")
	}

	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}
