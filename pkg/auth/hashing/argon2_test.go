package hashing

import (
	"regexp"
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	t.Parallel()

	hashRegex, err := regexp.Compile(`^\$argon2id\$v=19\$m=65536,t=1,p=4\$[A-Za-z0-9+/]{22}\$[A-Za-z0-9+/]{43}$`)
	if err != nil {
		t.Fatalf("failed to compile hash regex: %v", err)
	}

	hash1, err := HashPassword("t3stp4ssw0rd", DefaultArgon2Params)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if !hashRegex.MatchString(hash1) {
		t.Errorf("hash %q is not formatted properly", hash1)
	}

	hash2, err := HashPassword("t3stp4ssw0rd", DefaultArgon2Params)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if !hashRegex.MatchString(hash2) {
		t.Errorf("hash %q is not formatted properly", hash2)
	}

	if strings.Compare(hash1, hash2) == 0 {
		t.Errorf("hashes %q and %q are not unique", hash1, hash2)
	}
}

func TestComparePasswordAndHash(t *testing.T) {
	t.Parallel()

	hash := "$argon2id$v=19$m=65536,t=1,p=4$c8FqPHMT83tyxE2v0xDAFw$s2qmbRoRRbfyLIVFUzRwzE7F8PLjchpLKaV7Wf7tHgk"

	match, err := ComparePasswordAndHash("t3stp4ssw0rd", hash)
	if err != nil {
		t.Fatalf("failed to compare password and hash: %v", err)
	}

	if !match {
		t.Error("correct password and hash do not match")
	}

	match, err = ComparePasswordAndHash("wr0ngt3stp4ssw0rd", hash)
	if err != nil {
		t.Fatalf("failed to compare password and hash: %v", err)
	}

	if match {
		t.Error("wrong password and hash match")
	}
}

func TestComparePasswordAndInvalidHash(t *testing.T) {
	t.Parallel()

	_, err := ComparePasswordAndHash("t3stp4ssw0rd", "$argon2i$v=19$m=65536,t=1,p=4$c8FqPHMT83tyxE2v0xDAFw$s2qmbRoRRbfyLIVFUzRwzE7F8PLjchpLKaV7Wf7tHgk")
	if err != ErrInvalidArgon2Hash {
		t.Errorf("invalid error returned, got %v, want %v", err, ErrInvalidArgon2Hash)
	}

	_, err = ComparePasswordAndHash("t3stp4ssw0rd", "$argon2id$v=19$m=65536,t=1,p=4$c8FqPHMT83tyxE2v0xDAFw")
	if err != ErrInvalidArgon2Hash {
		t.Errorf("invalid error returned, got %v, want %v", err, ErrInvalidArgon2Hash)
	}

	_, err = ComparePasswordAndHash("t3stp4ssw0rd", "$argon2id$v=20$m=65536,t=1,p=4$c8FqPHMT83tyxE2v0xDAFw$s2qmbRoRRbfyLIVFUzRwzE7F8PLjchpLKaV7Wf7tHgk")
	if err != ErrIncompatibleArgon2Version {
		t.Errorf("invalid error returned, got %v, want %v", err, ErrIncompatibleArgon2Version)
	}

	_, err = ComparePasswordAndHash("t3stp4ssw0rd", "$argon2id$v=a$m=65536,t=1,p=4$c8FqPHMT83tyxE2v0xDAFw$s2qmbRoRRbfyLIVFUzRwzE7F8PLjchpLKaV7Wf7tHgk")
	if err != ErrIncompatibleArgon2Version {
		t.Errorf("invalid error returned, got %v, want %v", err, ErrIncompatibleArgon2Version)
	}

	_, err = ComparePasswordAndHash("t3stp4ssw0rd", "$argon2id$v=19$m=a,t=1,p=4$c8FqPHMT83tyxE2v0xDAFw$s2qmbRoRRbfyLIVFUzRwzE7F8PLjchpLKaV7Wf7tHgk")
	if err != ErrInvalidArgon2Hash {
		t.Errorf("invalid error returned, got %v, want %v", err, ErrInvalidArgon2Hash)
	}

	_, err = ComparePasswordAndHash("t3stp4ssw0rd", "$argon2id$v=19$m=65536,t=a,p=4$c8FqPHMT83tyxE2v0xDAFw$s2qmbRoRRbfyLIVFUzRwzE7F8PLjchpLKaV7Wf7tHgk")
	if err != ErrInvalidArgon2Hash {
		t.Errorf("invalid error returned, got %v, want %v", err, ErrInvalidArgon2Hash)
	}

	_, err = ComparePasswordAndHash("t3stp4ssw0rd", "$argon2id$v=19$m=65536,t=1,p=a$c8FqPHMT83tyxE2v0xDAFw$s2qmbRoRRbfyLIVFUzRwzE7F8PLjchpLKaV7Wf7tHgk")
	if err != ErrInvalidArgon2Hash {
		t.Errorf("invalid error returned, got %v, want %v", err, ErrInvalidArgon2Hash)
	}

	_, err = ComparePasswordAndHash("t3stp4ssw0rd", "$argon2id$v=19$m=65536,t=1,p=4$ä$s2qmbRoRRbfyLIVFUzRwzE7F8PLjchpLKaV7Wf7tHgk")
	if err != ErrInvalidArgon2Hash {
		t.Errorf("invalid error returned, got %v, want %v", err, ErrInvalidArgon2Hash)
	}

	_, err = ComparePasswordAndHash("t3stp4ssw0rd", "$argon2id$v=19$m=65536,t=1,p=4$c8FqPHMT83tyxE2v0xDAFw$ä")
	if err != ErrInvalidArgon2Hash {
		t.Errorf("invalid error returned, got %v, want %v", err, ErrInvalidArgon2Hash)
	}
}
