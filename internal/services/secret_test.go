package services

import (
	"testing"

	"github.com/stretchr/testify/require"
)



func TestSecret(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		masterKey string
	}{
		{
			name:      "Check coding and encodig func",
			data:      "secret data",
			masterKey: "Master",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bdata := []byte(tt.data)

			eKey, _, err := CreateEphemeralKey()
			require.NoError(t, err)

			// Check encode and decode eKey.
			eKeyc, err := EncodeKey(eKey, []byte(tt.masterKey))
			require.NoError(t, err)

			decodedEKey, err := DecodeKey(eKeyc, []byte(tt.masterKey))
			require.NoError(t, err)
			require.Equal(t, eKey, decodedEKey)

			// Check encode and decode dKey (data key).
			dKey, _, err := CreateDataKey()
			require.NoError(t, err)

			dKeyc, err := EncodeKey(dKey, decodedEKey)
			require.NoError(t, err)

			decodeddKey, err := DecodeKey(dKeyc, decodedEKey)
			require.NoError(t, err)
			require.Equal(t, dKey, decodeddKey)

			// Check encode and decode data.
			coded, err := EncodeData(decodeddKey, bdata)
			require.NoError(t, err)
			decoded, err := DecodeData(decodeddKey, coded)
			require.NoError(t, err)
			require.Equal(t, tt.data, string(decoded))
			t.Log("Data test: ", tt.data, " Decoded data: ", string(decoded))
		})
	}
}
