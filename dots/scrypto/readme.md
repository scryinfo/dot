# 说明

ed25519用于签名与验签  
x25519用于ecdh  
ed25519的private key与x25519的可以相互转换  
curve25519是一条曲线

# 使用x25519实现加解密的过程

1. 随机生成x25519的公私钥对A（每一次都会新生成）
2. 使用A的私钥及B的公钥计算出key
3. key通过HKDF混合得到wrapped key
4. 使用AEC算法用wrpped key加密数据

5. B使用A的公钥计算出key
6. key通过HKDF混合得到wrapped key
7. 使用AEC算法解密数据

# 参见

[age](https://github.com/FiloSottile/age)
[ecdh25519](https://github.com/aead/ecdh)