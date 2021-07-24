package sasl

import (
	"crypto/hmac"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
	"strings"
)

// Implementation of SCRAM (RFC 5802)
type Scram struct {
	// gs2Header string
	username string
	nonce    string
	proof    []byte // sent from client

	// used for computing serverSignature
	clientFirstBare, serverFirst, clientFinalWithoutProof string

	// the Cred associated with the requesting client
	Cred *Credential

	// Hash function (`H()` in RFC 5802)
	Hash func() hash.Hash
}

func (s *Scram) ParseClientFirst(m string) error {
	attrs := strings.Split(m, ",")
	if len(attrs) < 4 {
		return errors.New("e=other-error")
	}

	// attrs[1] is unused as we do not take advantage of authzid

	s.username = attrs[2][2:]

	// add arbitrary length nonce
	nonce := make([]byte, 20)
	rand.Read(nonce)
	s.nonce = attrs[3][2:] + base64.StdEncoding.EncodeToString(nonce)

	s.clientFirstBare = strings.Join(attrs[2:], ",")
	return nil
}

func (s *Scram) GenServerFirst() string {
	s.serverFirst = fmt.Sprintf("r=%s,s=%s,i=%d",
		s.nonce,
		base64.StdEncoding.EncodeToString(s.Cred.Salt),
		s.Cred.Iteration,
	)
	return s.serverFirst
}

func (s *Scram) ParseClientFinal(m string) error {
	attrs := strings.Split(m, ",")
	if len(attrs) < 3 {
		return errors.New("e=other-error")
	}

	// attrs[0] is unused since we don't use channel binding
	nonce := attrs[1][2:]
	if nonce != s.nonce {
		return errors.New("e=other-error")
	}

	proof, err := base64.StdEncoding.DecodeString(attrs[2][2:])
	if err != nil {
		return errors.New("e=invalid-encoding")
	}
	s.proof = proof

	s.clientFinalWithoutProof = strings.Join(attrs[:2], ",")
	return nil
}

func (s *Scram) GenServerFinal() (string, error) {
	authMsg := fmt.Sprintf("%s,%s,%s", s.clientFirstBare, s.serverFirst, s.clientFinalWithoutProof)

	mac := hmac.New(s.Hash, s.Cred.StoredKey)
	mac.Write([]byte(authMsg))
	clientSignature := mac.Sum(nil)

	clientKey := bytewiseXOR(clientSignature, s.proof)

	hash := s.Hash()
	hash.Write(clientKey)
	storedKey := hash.Sum(nil)

	if !hmac.Equal(storedKey, s.Cred.StoredKey) {
		return "", errors.New("e=invalid-proof")
	}

	mac = hmac.New(s.Hash, s.Cred.ServerKey)
	mac.Write([]byte(authMsg))
	serverSignature := mac.Sum(nil)

	return fmt.Sprintf("v=%s", base64.StdEncoding.EncodeToString(serverSignature)), nil
}

func bytewiseXOR(b1, b2 []byte) []byte {
	if len(b1) != len(b2) {
		return nil
	}

	x := make([]byte, len(b1))
	for i := range x {
		x[i] = b1[i] ^ b2[i]
	}
	return x
}

func SCRAM(c *Credential, h func() hash.Hash) *Scram {
	return &Scram{Cred: c, Hash: h}
}
