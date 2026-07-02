package delivery

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/mc-lovin-132/auth/internal/domain"

	"google.golang.org/grpc/metadata"
)

type DeviceDetector struct {
	salt string
}

func New(salt string) *DeviceDetector {
	return &DeviceDetector{salt: salt}
}

const stringZeroValue = ""

// для определения устройства нам не нужны никакие данные, все обходимое есть в контексте
// чтобы при обращении через апи gateway все работало, нужно пробрасывать данные
func (d *DeviceDetector) Detect(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return stringZeroValue, fmt.Errorf("%w:%s", domain.ErrInvalidRequest, "no metadata in request")
	}
	userAgent, err := getFirst(md, "user-agent")
	if err != nil {
		return stringZeroValue, err
	}
	data := []byte(userAgent + " " + d.salt)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

func getFirst(md metadata.MD, key string) (string, error) {
	values := md.Get(key)
	if len(values) == 0 {
		return stringZeroValue, fmt.Errorf("%w:%s", domain.ErrInvalidRequest, "no metadata in request")
	}
	return values[0], nil
}
