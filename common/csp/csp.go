package csp

import (
	"encoding/hex"
	"log"

	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/bccsp/factory"
)

var prov bccsp.BCCSP

func init() {
	var err error
	if prov, err = NewCSP("."); err != nil {
		log.Fatalln(err)
	}
}

// DefaultCSP return default provider
func DefaultCSP() bccsp.BCCSP {
	return prov
}

// NewCSP create a crpto service provider
func NewCSP(path string) (csp bccsp.BCCSP, err error) {
	opts := &factory.FactoryOpts{
		ProviderName: "SW",
		SwOpts: &factory.SwOpts{
			HashFamily: "SHA2",
			SecLevel:   256,

			FileKeystore: &factory.FileKeystoreOpts{
				KeyStorePath: path,
			},
		},
	}

	csp, err = factory.GetBCCSPFromOpts(opts)
	return
}

func generatePrivKey(temp bool) (bccsp.Key, error) {
	return prov.KeyGen(&bccsp.ECDSAP256KeyGenOpts{Temporary: temp})
}

// GeneratePrivateKey generate private key
func GeneratePrivateKey() (bccsp.Key, error) {
	return generatePrivKey(false)
}

// GeneratePrivateTempKey generate a temporary private key
func GeneratePrivateTempKey() (bccsp.Key, error) {
	return generatePrivKey(true)
}

//func GeneratePrivateKey(csp bccsp.BCCSP) (bccsp.Key, error) {
//	return csp.KeyGen(&bccsp.ECDSAP256KeyGenOpts{Temporary: false})
//}

// Hash SHA256 Hash
func Hash(msg []byte) ([]byte, error) {
	return prov.Hash(msg, &bccsp.SHA256Opts{})
}

//func Hash(prov bccsp.BCCSP, msg []byte) ([]byte, error) {
//	return prov.Hash(msg, &bccsp.SHA256Opts{})
//}

// Sign sign a message
func Sign(priv bccsp.Key, digest []byte) ([]byte, error) {
	return prov.Sign(priv, digest, nil)
}

//func Sign(prov bccsp.BCCSP, priv bccsp.Key, digest []byte) ([]byte, error) {
//	return prov.Sign(priv, digest, nil)
//}

// Verify verify a signature is valid or not
func Verify(key bccsp.Key, signature, digest []byte) (bool, error) {
	return prov.Verify(key, signature, digest, nil)
}

//func Verify(prov bccsp.BCCSP, priv bccsp.Key, signature, digest []byte) (bool, error) {
//	return prov.Verify(priv, signature, digest, nil)
//}

// ImportPrivFromSKI get a private key from SKI string
func ImportPrivFromSKI(ski string) (bccsp.Key, error) {
	bytes, err := hex.DecodeString(ski)

	if err == nil {
		return prov.GetKey(bytes)
	}

	return nil, err
}

// ImportPubFromRaw import a public key from raw bytes
func ImportPubFromRaw(raw []byte) (bccsp.Key, error) {
	return prov.KeyImport(raw, &bccsp.ECDSAPKIXPublicKeyImportOpts{Temporary: true})
}

// GetPublicKeyBytes get public key bytes from private key
func GetPublicKeyBytes(priv bccsp.Key) ([]byte, error) {
	pub, err := priv.PublicKey()
	if err != nil {
		return nil, err
	}

	return pub.Bytes()
}
