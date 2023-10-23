package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// PlayerInfoKeyPrefix is the prefix to retrieve all PlayerInfo
	PlayerInfoKeyPrefix = "PlayerInfo/value/"
)

// PlayerInfoKey returns the store key to retrieve a PlayerInfo from the index fields
func PlayerInfoKey(
	account string,
) []byte {
	var key []byte

	accountBytes := []byte(account)
	key = append(key, accountBytes...)
	key = append(key, []byte("/")...)

	return key
}
