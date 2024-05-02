package services

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"slices"
	"time"

	"github.com/gofrs/uuid"
	"github.com/shulganew/GophKeeper/internal/entities"
	"go.uber.org/zap"
)

const DataKeyLen = 8
const EphemeralKeyLen = 16

// Return actual currently using key for encoding data keys.
func (k *Keeper) GetActualEKey() (eKey entities.EKeyMem) {
	eKey = slices.MaxFunc(k.eKeys, func(a, b entities.EKeyMem) int {
		if a.TS.After(b.TS) {
			return 1
		}
		return 0
	})
	return
}

// Get key from memory by time stamp (ts used as key version).
func (k *Keeper) GetEKey(ts time.Time) (eKey entities.EKeyMem) {
	id := slices.IndexFunc(k.eKeys, func(key entities.EKeyMem) bool { return key.TS.Equal(ts) })
	if id != -1 {
		return k.eKeys[id]
	}
	return
}

// Load all keys to ekeys key ring, (last time stamp), with master key encoding.
func (k *Keeper) LoadKeyRing(ctx context.Context) {
	eKeysc, err := k.stor.LoadEKeysc(ctx)
	if err != nil {
		zap.S().Errorln("Error load keys: ", err)
	}
	// No created eKeys.
	if len(eKeysc) == 0 {
		eKey, ts, err := CreateEphemeralKey()
		if err != nil {
			zap.S().Errorln("Error create  eKeys: ", err)
		}
		// Load to mem.
		ekeyMem := entities.EKeyMem{EKey: eKey, TS: ts}
		k.eKeys = []entities.EKeyMem{ekeyMem}
		// Save to storage.
		k.SaveKeyRing(ctx, ekeyMem)
		return
	}

	// Decode eKeyc to eKey.
	eKeys := []entities.EKeyMem{}
	for _, eKeyc := range eKeysc {
		// Decode keys using master key.
		eKey, err := DecodeKey(eKeyc.EKeyc, []byte(k.conf.MasterKey))
		if err != nil {
			zap.S().Errorln("Error Decode eKeys: ", err)
		}
		eKeys = append(eKeys, entities.EKeyMem{TS: eKeyc.TS, EKey: eKey})
	}
	// Init memory.
	k.eKeys = eKeys
}

// Save key from ekey  to database with master key encoding.
func (k *Keeper) SaveKeyRing(ctx context.Context, eKey entities.EKeyMem) {

	// Encode keys using master key.
	eKeyc, err := EncodeKey(eKey.EKey, []byte(k.conf.MasterKey))
	if err != nil {
		zap.S().Errorln("Error Encoding eKey: ", err)
	}
	// Save to database coded key.
	err = k.stor.SaveEKeyc(ctx, entities.EKeyDB{TS: eKey.TS, EKeyc: eKeyc})
	if err != nil {
		zap.S().Errorln("Error save eKeyc to DB: ", err)
	}
}

// Create new Ephemeral key.
func CreateEphemeralKey() (eKey []byte, ts time.Time, err error) {
	eKey, err = createKey(EphemeralKeyLen)
	ts = time.Now()
	return
}

// Create Data key. Key saved in data table "secrets", column "data_key"
func CreateDataKey() (dKey []byte, ts time.Time, err error) {
	dKey, err = createKey(DataKeyLen)
	ts = time.Now()
	return
}

// Create key particular size.
func createKey(size int) (key []byte, err error) {
	data := make([]byte, size)
	_, err = rand.Read(data)
	if err != nil {
		return nil, err
	}
	return data, nil

}

// Encode original key with key, i.e. eKey encoded with master key from memory to store.
func EncodeKey(storingKey, useKey []byte) ([]byte, error) {
	coded, err := EncodeData(useKey, []byte(storingKey))
	if err != nil {
		return nil, err
	}
	return coded, nil
}

// Dekode original key with key, i.e. eKey decoded with master key from store to memory.
func DecodeKey(storingKey, useKey []byte) ([]byte, error) {
	coded, err := DecodeData(useKey, []byte(storingKey))
	if err != nil {
		return nil, err
	}
	return coded, nil
}

// Encode data using AES with string key.
func EncodeData(dKey []byte, data []byte) (coded []byte, err error) {
	nonce, aesgcm, err := getCryptData(dKey)
	if err != nil {
		zap.S().Errorln("Encription Error: get enctypt data")
		return
	}
	coded = aesgcm.Seal(nil, nonce, data, nil)

	return
}

// Decode data using AES with string key.
func DecodeData(dKey []byte, coded []byte) (data []byte, err error) {
	nonce, aesgcm, err := getCryptData(dKey)
	if err != nil {
		zap.S().Errorln("Encription Error: get enctypt data")
		return
	}

	data, err = aesgcm.Open(nil, nonce, coded, nil)
	if err != nil {
		zap.S().Errorln("Encryption Error: Open seal")
		return
	}
	return
}

// Get nonce and cipher from string, help func
func getCryptData(key []byte) (nonce []byte, aesgcm cipher.AEAD, err error) {
	hkey := sha256.Sum256(key)

	aesblock, err := aes.NewCipher(hkey[:32])
	if err != nil {
		zap.S().Errorln("Encryption Error: aesblock")
		return
	}

	aesgcm, err = cipher.NewGCM(aesblock)
	if err != nil {
		zap.S().Errorln("Encryption Error: aesgcm")
		return
	}

	length := aesgcm.NonceSize()
	nonceSize := len(hkey) - length
	nonce = hkey[nonceSize:]
	return
}

// Common method for all data types to store cypted data in DB.
func (k *Keeper) AddSecret(ctx context.Context, userID string, dataType entities.SecretType, data []byte) (secretID *uuid.UUID, err error) {
	// Get data key.
	dKey, _, err := CreateDataKey()
	if err != nil {
		zap.S().Errorln("Error create data key: ", err)
		return nil, err
	}
	// Encode date before store.
	datac, err := EncodeData(dKey, data)
	if err != nil {
		zap.S().Errorln("Error encode data: ", err)
		return nil, err
	}
	// Get Ephemeral current key.
	eKey := k.GetActualEKey()
	dKeyc, err := EncodeKey(dKey, eKey.EKey)
	if err != nil {
		zap.S().Errorln("Error getting ephemeral key: ", err)
		return nil, err
	}

	dbSite := entities.NewSecretEncoded{NewSecret: entities.NewSecret{UserID: userID, Type: dataType, EKeyVer: eKey.TS, DKey: dKeyc, Uploaded: time.Now()}, DataCr: datac}

	secretID, err = k.stor.AddSecretStor(ctx, dbSite, dataType)
	if err != nil {
		zap.S().Errorln("Error adding site credentials: ", err)
		return nil, err
	}
	return
}

func (k *Keeper) GetSecrets(ctx context.Context, userID string, dataType entities.SecretType) (secrets []entities.SecretDecoded, err error) {
	// Load all user's sites coded credentials from database.
	secretsc, err := k.stor.GetSecretStor(ctx, userID, dataType)
	if err != nil {
		zap.S().Errorln("Error getting site credentials: ", err)
		return nil, err
	}

	for _, secret := range secretsc {
		// Get ephemeral key (version from ts in db) for decode data key.
		eKey := k.GetEKey(secret.EKeyVer)
		// Decode dKeyc
		dKey, err := DecodeKey(secret.DKey, eKey.EKey)
		if err != nil {
			zap.S().Errorln("Error decode data key: ", err)
			return nil, err
		}
		// Decode data.
		data, err := DecodeData(dKey, secret.DataCr)
		if err != nil {
			zap.S().Errorln("Error decode stored data: ", err)
			return nil, err
		}

		secretDecoded := entities.SecretDecoded{NewSecret: entities.NewSecret{UserID: userID, Type: dataType, EKeyVer: secret.EKeyVer, Uploaded: secret.Uploaded}, UUID: secret.UUID, Data: data}
		//Definition: newSite.Definition, SiteID: dbSite.UUID.String(), Site: newSite.Site, Slogin: newSite.Slogin, Spw: newSite.Spw}
		secrets = append(secrets, secretDecoded)
	}
	return
}
