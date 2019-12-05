package factory

import (
	"io"

	"github.com/golang/glog"
	"golang.org/x/crypto/openpgp"
)

// BuildNullDecrypt creates a method which fits the profile of a tetra
// decryption method, but does not make any change to the data passed.
func BuildNullDecrypt() func(io.Reader) (io.Reader, error) {
	nullDecrypt := func(encrypted io.Reader) (io.Reader, error) {
		return encrypted, nil
	}

	return nullDecrypt
}

// BuildPGPDecrypt creates a method which returns a PGP decrypted io.Reader
// based on the provided KeyRing.
func BuildPGPDecrypt(
	privateKey io.Reader,
	passphrase []byte,
) func(io.Reader) (io.Reader, error) {

	pgpDecrypt := func(encrypted io.Reader) (io.Reader, error) {
		var entity *openpgp.Entity
		var entityList openpgp.EntityList
		entityList, err := openpgp.ReadArmoredKeyRing(privateKey)
		if err != nil {
			glog.Fatalf("error: could not read key ring (%s)", err)
		}

		entity = entityList[0]

		err = entity.PrivateKey.Decrypt(passphrase)
		if err != nil {
			glog.Fatalf("error: could not decrypt private key (%s)", err)
		}
		for _, subkey := range entity.Subkeys {
			err = subkey.PrivateKey.Decrypt(passphrase)
			if err != nil {
				glog.Fatalf("error: could not decrypt private subkey (%s)", err)
			}
		}
		glog.Infof("encrypted: %+v \n%+v\n", encrypted, *entity)
		md, err := openpgp.ReadMessage(encrypted, entityList, nil, nil)
		if err != nil {
			glog.Errorf("error: could not decrypt the file using pgp (%+v %s)", md, err)
		}

		return md.UnverifiedBody, nil
	}

	return pgpDecrypt
}
