package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"time"
	"errors"
	mathRand "math/rand"
	"strings"
	"bytes"
	"compress/gzip"
	"math/big"
	"crypto/x509"
	"encoding/pem"
	"encoding/base64"
)

type XEcdsa struct {
	publicKey  *ecdsa.PublicKey
	privateKey *ecdsa.PrivateKey
}

func NewXEcdsa(publicKey []byte, privateKey []byte) (*XEcdsa, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("block is error")
	}

	priKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	block, _ = pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("block is error")
	}

	pubInner, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pubKey := pubInner.(*ecdsa.PublicKey)

	return &XEcdsa{
		privateKey: priKey,
		publicKey:  pubKey,
	}, nil
}

//生成指定math/rand字节长度的随机字符串
func GetRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*()_+?=-"
	bytes := []byte(str)
	result := []byte{}
	r := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//生成ECC算法的公钥和私钥文件
//根据随机字符串生成，randKey至少36位
func GenerateKey(randKey string) ([]byte, []byte, error) {

	var err error
	var privateKey *ecdsa.PrivateKey
	var publicKey *ecdsa.PublicKey
	var curve elliptic.Curve

	//一、生成私钥文件
	//根据随机字符串长度设置curve曲线
	length := len(randKey)
	//elliptic包实现了几条覆盖素数有限域的标准椭圆曲线,Curve代表一个短格式的Weierstrass椭圆曲线，其中a=-3
	if length < 224/8 {
		err = errors.New("私钥长度太短，至少为36位！")
		return nil, nil, err
	}

	if length >= 521/8+8 {
		//长度大于73字节，返回一个实现了P-512的曲线
		curve = elliptic.P521()
	} else if length >= 384/8+8 {
		//长度大于56字节，返回一个实现了P-384的曲线
		curve = elliptic.P384()
	} else if length >= 256/8+8 {
		//长度大于40字节，返回一个实现了P-256的曲线
		curve = elliptic.P256()
	} else if length >= 224/8+8 {
		//长度大于36字节，返回一个实现了P-224的曲线
		curve = elliptic.P224()
	}

	//GenerateKey方法生成私钥
	privateKey, err = ecdsa.GenerateKey(curve, strings.NewReader(randKey))
	if err != nil {
		return nil, nil, err
	}

	//通过x509标准将得到的ecc公钥序列化为ASN.1的DER编码字符串
	priKey, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, nil, nil
	}

	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: priKey,
	}
	pri := pem.EncodeToMemory(block)

	//二、生成公钥文件
	//从得到的私钥对象中将公钥信息取出
	publicKey = &privateKey.PublicKey
	//通过x509标准将得到的ecc公钥序列化为ASN.1的DER编码字符串
	pubKey, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, nil, nil
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKey,
	}
	pub := pem.EncodeToMemory(block)

	return pri, pub, nil
}

//使用ECC算法加密签名，返回签名数据
func (e *XEcdsa) CryptSignByEcc(input string, randSign string) (output string, err error) {
	//ecc私钥和随机签字符串数据得到哈希
	r, s, err := ecdsa.Sign(strings.NewReader(randSign), e.privateKey, []byte(input))
	if err != nil {
		return "", err
	}

	rt, err := r.MarshalText()
	if err != nil {
		return "", err
	}

	st, err := s.MarshalText()
	if err != nil {
		return "", err
	}

	//拼接两个椭圆曲线参数哈希
	var b bytes.Buffer
	writer := gzip.NewWriter(&b)
	defer writer.Close()

	_, err = writer.Write([]byte(string(rt) + "+" + string(st)))
	if err != nil {
		return "", err
	}
	writer.Flush()

	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

//使用ECC算法,对密文和明文进行匹配校验
func (e *XEcdsa) VerifyCryptEcc(srcStr, cryptStr string) (bool, error) {
	decode, err := base64.StdEncoding.DecodeString(cryptStr)
	if err != nil {
		return false, err
	}
	//解密签名信息，返回椭圆曲线参数：两个大整数
	rint, sint, err := UnSignCryptEcc(decode)
	//使用公钥、原文、以及签名信息解密后的两个椭圆曲线的大整数参数进行校验
	verify := ecdsa.Verify(e.publicKey, []byte(srcStr), &rint, &sint)

	return verify, nil
}

//使用ECC算法解密,返回加密前的椭圆曲线大整数
func UnSignCryptEcc(cryptBytes []byte) (rint, sint big.Int, err error) {
	reader, err := gzip.NewReader(bytes.NewBuffer(cryptBytes))
	if err != nil {
		err = errors.New("decode error," + err.Error())
	}
	defer reader.Close()

	buf := make([]byte, 1024)
	count, err := reader.Read(buf)
	if err != nil {
		err = errors.New("decode read error," + err.Error())
	}

	rs := strings.Split(string(buf[:count]), "+")
	if len(rs) != 2 {
		err = errors.New("decode fail")
		return
	}
	err = rint.UnmarshalText([]byte(rs[0]))
	if err != nil {
		err = errors.New("decrypt rint fail, " + err.Error())
		return
	}
	err = sint.UnmarshalText([]byte(rs[1]))
	if err != nil {
		err = errors.New("decrypt sint fail, " + err.Error())
		return
	}
	return
}
