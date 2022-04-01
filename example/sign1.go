package main

import (
	"crypto/rand"
	_ "crypto/sha256"
	"log"

	cose "github.com/veraison/go-cose"
)

func main() {
	msg := cose.NewSign1Message()

	// EAT token from Appendix A.2 of draft-ietf-rats-eat
	msg.Payload = []byte{
		0xa7, 0x0a, 0x49, 0x94, 0x8f, 0x88, 0x60, 0xd1, 0x3a, 0x46,
		0x3e, 0x8e, 0x0b, 0x53, 0x01, 0x98, 0xf5, 0x0a, 0x4f, 0xf6,
		0xc0, 0x58, 0x61, 0xc8, 0x86, 0x0d, 0x13, 0xa6, 0x38, 0xea,
		0x4f, 0xe2, 0xfa, 0x0f, 0xf5, 0x10, 0x03, 0x06, 0xc1, 0x1a,
		0x5a, 0xfd, 0x32, 0x2e, 0x0e, 0x03, 0x14, 0xa3, 0x6f, 0x41,
		0x6e, 0x64, 0x72, 0x6f, 0x69, 0x64, 0x20, 0x41, 0x70, 0x70,
		0x20, 0x46, 0x6f, 0x6f, 0xa1, 0x0e, 0x01, 0x72, 0x53, 0x65,
		0x63, 0x75, 0x72, 0x65, 0x20, 0x45, 0x6c, 0x65, 0x6d, 0x65,
		0x6e, 0x74, 0x20, 0x45, 0x61, 0x74, 0xd8, 0x3d, 0xd2, 0x42,
		0x01, 0x23, 0x6d, 0x4c, 0x69, 0x6e, 0x75, 0x78, 0x20, 0x41,
		0x6e, 0x64, 0x72, 0x6f, 0x69, 0x64, 0xa1, 0x0e, 0x01,
	}

	// create a signer with a new private key
	// ES256 (algId: -7), i.e.: ECDSA w/ SHA-256 from RFC8152
	signer, err := cose.NewSigner(cose.ES256, nil)
	if err != nil {
		log.Fatalf("signer creation failed: %s", err)
	}

	msg.Headers.Protected[1] = -7 // ECDSA w/ SHA-256

	msg.Headers.Unprotected["kid"] = 1

	// no external data
	err = msg.Sign(rand.Reader, nil, *signer)
	if err != nil {
		log.Fatalf("signature creation failed: %s\n", err)
	}

	log.Printf("COSE Sign1 signature bytes: %x\n", msg.Signature)

	coseSig, err := cose.Marshal(msg)
	if err != nil {
		log.Fatalf("COSE marshaling failed: %s", err)
	}

	log.Printf("COSE Sign1 message: %x\n", coseSig)

	// derive a verifier using the signer's public key and COSE algorithm
	verifier := signer.Verifier()

	// Verify
	err = msg.Verify(nil, *verifier)
	if err != nil {
		log.Fatalf("Error verifying the message %+v", err)
	}

	log.Println("COSE Sign1 signature verified")
}