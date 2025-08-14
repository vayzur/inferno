package auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/vayzur/inferno/internal/config"
)

func VerifyRollingHash(header string) error {
	if !strings.HasPrefix(header, "rolling ") {
		return errors.New("invalid header prefix")
	}

	auth := strings.TrimPrefix(header, "rolling ")
	parts := strings.SplitN(auth, ":", 2)
	if len(parts) != 2 {
		return errors.New("invalid format")
	}

	tsStr, sig := parts[0], parts[1]
	ts, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		return errors.New("invalid timestamp")
	}

	if time.Since(time.Unix(ts, 0)).Abs() > time.Minute {
		return errors.New("expired")
	}

	var b []byte
	b = fmt.Appendf(b, "%d:%s", ts, config.AppConfig.Token)
	hash := sha256.Sum256(b)
	expected := hex.EncodeToString(hash[:])

	if subtle.ConstantTimeCompare([]byte(sig), []byte(expected)) == 1 {
		return nil
	}

	return errors.New("unauthorized")
}
