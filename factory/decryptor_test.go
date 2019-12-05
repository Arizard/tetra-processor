package factory

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"golang.org/x/crypto/openpgp/armor"
)

func TestBuildPGPDecrypt(t *testing.T) {
	type args struct {
		privateKey io.Reader
		passphrase []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"Test decryption using pgp",
			args{
				strings.NewReader(privateKeyString),
				passphrase,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pgpDecrypt := BuildPGPDecrypt(tt.args.privateKey, tt.args.passphrase)

			encryptedDecoded, err := armor.Decode(encrypted)
			if err != nil {
				t.Errorf("error: pgp armor was not decoded (%s)", err)
			}

			decrypted, err := pgpDecrypt(encryptedDecoded.Body)
			if err != nil {
				t.Errorf("error: pgpDecrypt failed (%s)", err)
			}

			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(decrypted)
			if err != nil {
				t.Errorf("error: buffer could not read decrypted (%s)", err)
			}
			if buf.String() != decryptedString {
				t.Errorf("error: pgpDecrypt had incorrect output (got %s, want %s)", buf.String(), decryptedString)
			}
			// if got := BuildPGPDecrypt(tt.args.privateKey, tt.args.passphrase); !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("BuildPGPDecrypt() = %v, want %v", got, tt.want)
			// }
		})
	}
}

const privateKeyString = `-----BEGIN PGP PRIVATE KEY BLOCK-----

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
=XUIo
-----END PGP PRIVATE KEY BLOCK-----
`
const encryptedString = `-----BEGIN PGP MESSAGE-----

hQEMA7+OPh9ZnJ8/AQf8DR5r1lKh1UTHMpSTPL4D+xA7tDmV9p1ehl3MU0Uht4Hf
4Tti8o2k+I2sx30M5kAWIE/4SSYSqsH7Ztl4giElGIDat+NP+sfGEPCNd4vZO+K+
dYlLJW7fzqvgVY1WhP/m6qCVRMD0WfdYq4HRnJZTxYP1b7x+nk5Z0a0QuomyQvCT
fRmw8Njtx6UwgCuhr0I7UTBhFn2EIYzJDMR7vPjBbSmImeXNNUGQ/E7yZuuEzLiL
78RanL+pEhbMEhUW0PfVTwb0qCqwGRTVuR1toO0eejgI5cS1YZZhbF1vp2MmjKXh
bXNlNvHKPMF0WK0HylyzsHPH7VQsGhoajZGqcdDAftJWARTUKplMqUnjvisP2UJr
i4ceDSY+uvC/1ZVm8+oteCJukBOAzr6vj7k3pFhObWD2crVTVl0L71PE8QQAw2TX
OmgyPfdgrNfUH/7c9SN4rl4B2xlMGec=
=8su1
-----END PGP MESSAGE-----
`
const decryptedString = `hello mr tetra
`

var passphrase = []byte("W0lfSp1d3r")
var encrypted = strings.NewReader(encryptedString)
