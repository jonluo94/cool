package crypto

import (
	"encoding/base64"
	"crypto/x509"
	"crypto/rsa"
	"crypto/rand"
	"errors"
	"crypto"
	"bytes"
	"encoding/asn1"
	"encoding/pem"
)

const (
	RSA_ALGORITHM_SIGN = crypto.SHA256
)

type XRsa struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

// 生成密钥对
func CreateKeys() (pri []byte, pub []byte, err error) {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, nil, err
	}

	derPriKey := MarshalPKCS8PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derPriKey,
	}
	pri = pem.EncodeToMemory(block)

	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPubKey, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, nil, err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPubKey,
	}
	pub = pem.EncodeToMemory(block)
	return
}
func NewXRsa(publicKey []byte, privateKey []byte) (*XRsa, error) {

	//加载公钥
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("block is error")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubKey.(*rsa.PublicKey)

	//解析私钥
	block, _ = pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("block is error")
	}

	priKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pri, ok := priKey.(*rsa.PrivateKey)

	if ok {
		return &XRsa{
			publicKey:  pub,
			privateKey: pri,
		}, nil
	} else {
		return nil, errors.New("private key not supported")
	}
}

// 公钥加密
func (r *XRsa) PublicEncrypt(data string) (string, error) {
	partLen := r.publicKey.N.BitLen()/8 - 11
	chunks := split([]byte(data), partLen)
	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		bytes, err := rsa.EncryptPKCS1v15(rand.Reader, r.publicKey, chunk)
		if err != nil {
			return "", err
		}
		buffer.Write(bytes)
	}
	return base64.RawURLEncoding.EncodeToString(buffer.Bytes()), nil
}

// 私钥解密
func (r *XRsa) PrivateDecrypt(encrypted string) (string, error) {
	partLen := r.publicKey.N.BitLen() / 8
	raw, err := base64.RawURLEncoding.DecodeString(encrypted)
	chunks := split([]byte(raw), partLen)
	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, r.privateKey, chunk)
		if err != nil {
			return "", err
		}
		buffer.Write(decrypted)
	}
	return buffer.String(), err
}

// 数据加签
func (r *XRsa) Sign(data string) (string, error) {
	h := RSA_ALGORITHM_SIGN.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)
	sign, err := rsa.SignPKCS1v15(rand.Reader, r.privateKey, RSA_ALGORITHM_SIGN, hashed)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(sign), err
}

// 数据验签
func (r *XRsa) Verify(data string, sign string) error {
	h := RSA_ALGORITHM_SIGN.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)
	decodedSign, err := base64.RawURLEncoding.DecodeString(sign)
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(r.publicKey, RSA_ALGORITHM_SIGN, hashed, decodedSign)
}
func MarshalPKCS8PrivateKey(key *rsa.PrivateKey) []byte {
	info := struct {
		Version             int
		PrivateKeyAlgorithm []asn1.ObjectIdentifier
		PrivateKey          []byte
	}{}
	info.Version = 0
	info.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
	info.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	info.PrivateKey = x509.MarshalPKCS1PrivateKey(key)
	k, _ := asn1.Marshal(info)
	return k
}
func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}
