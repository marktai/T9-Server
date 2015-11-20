package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"db"
	"encoding/base64"
	"encoding/hex"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"math/big"
	mrand "math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var secretMap = make(map[uint]*Secret)

type Secret struct {
	value      *big.Int
	expiration time.Time
}

// secrets are an arbitrary big int number from 0 to 2^512
// to actually use their value, they are converted into base64
// then the base64 string chararcters are used as bytes
// this is to get random bytes and still be able to nicely store them in strings
func (s *Secret) Bytes() []byte {
	return s.value.Bytes()
}

func (s *Secret) Base64() string {
	return base64.StdEncoding.EncodeToString(s.Bytes())
}

func (s *Secret) String() string {
	return s.Base64()
}

var bitSize int64 = 512

var limit *big.Int

func newSecret() (*Secret, error) {
	if limit == nil {
		limit = big.NewInt(2)
		limit.Exp(big.NewInt(2), big.NewInt(bitSize), nil)
	}

	value, err := rand.Int(rand.Reader, limit)
	if err != nil {
		return nil, err
	}
	retSecret := &Secret{value, time.Now()}
	return retSecret, nil
}

func stringtoUint(s string) (uint, error) {
	i, err := strconv.Atoi(s)
	return uint(i), err
}

func getUniqueID() (uint, error) {

	mrand.Seed(time.Now().Unix())
	collision := 1
	times := 0
	var id uint

	for collision != 0 {

		id = uint(mrand.Int31n(65536))
		err := db.Db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id=?)", id).Scan(&collision)
		if err != nil {
			return id, err
		}
		times++
		if times > 20 {
			return id, errors.New("Too many attempts to find a unique game ID")
		}
	}
	return id, nil
}

func MakeUser(user, pass string) (uint, error) {
	_, err := getUserID(user)
	if err == nil {
		return 0, errors.New("User already made")
	}

	id, err := getUniqueID()
	if err != nil {
		return 0, err
	}

	saltHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	saltHashString := base64.StdEncoding.EncodeToString(saltHash)

	err = db.Db.Ping()
	if err != nil {
		return 0, err
	}

	addUser, err := db.Db.Prepare("INSERT INTO users VALUES(?, ?, ?)")
	if err != nil {
		return 0, err
	}

	_, err = addUser.Exec(id, user, saltHashString)

	if err != nil {
		return 0, err
	}

	return id, nil

}

func getUserID(user string) (uint, error) {

	var userID uint
	err := db.Db.QueryRow("SELECT id FROM users WHERE name=?", user).Scan(&userID)

	return userID, err
}

func getSaltHash(userID uint) ([]byte, error) {
	saltHashString := ""
	err := db.Db.QueryRow("SELECT salthash FROM users WHERE id=?", userID).Scan(&saltHashString)
	if err != nil {
		return nil, err
	}
	saltHash, err := base64.StdEncoding.DecodeString(saltHashString)
	return saltHash, err
}

func Login(user, pass string) (uint, *Secret, error) {

	userID, err := getUserID(user)
	hash, err := getSaltHash(userID)
	if err != nil {
		return 0, nil, err
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(pass))

	if err != nil {
		return 0, nil, err
	}

	if _, ok := secretMap[userID]; !ok {
		secret, err := newSecret()
		if err != nil {
			return 0, nil, err
		}
		secretMap[userID] = secret
	}

	return userID, secretMap[userID], nil
}

func ComputeHmac256(message, key string) []byte {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	expectedMAC := mac.Sum(nil)

	return expectedMAC
}

// CheckMAC reports whether messageHMAC is a valid HMAC tag for message.
func checkMAC(message, key string, messageHMAC []byte) bool {
	expectedMAC := ComputeHmac256(message, key)
	return hmac.Equal(messageHMAC, expectedMAC)
}

// returns userID, message used to generate HMAC, and HMAC from request
func parseRequest(r *http.Request) (string, string, []byte, error) {
	timeSlice, ok := r.Header["Time-Sent"]
	if !ok || timeSlice == nil {
		return "", "", nil, errors.New("No Time-Sent header provided")
	}

	time := timeSlice[0]

	message := time + ":" + r.URL.String()

	messageHMACSlice, ok := r.Header["Hmac"]
	if !ok || timeSlice == nil {
		return "", "", nil, errors.New("No HMAC header provided")
	}

	messageHMACString := messageHMACSlice[0]
	HMACFormat := ""

	formatSlice, ok := r.Header["Format"]
	if ok && formatSlice != nil {
		format := formatSlice[0]
		if ok {
			switch strings.ToLower(format) {
			case "base64", "64":
				HMACFormat = "base64"
			case "hex", "hexadecimal":
				HMACFormat = "hex"
			case "binary", "bits":
				HMACFormat = "binary"
			case "decimal":
				HMACFormat = "decimal"
			default:
				HMACFormat = format
			}
		}
	} else {
		HMACFormat = "base64"
	}

	var messageHMAC []byte
	var err error
	switch HMACFormat {
	case "base64":
		messageHMAC, err = base64.StdEncoding.DecodeString(messageHMACString)
	case "hex":
		messageHMAC, err = hex.DecodeString(messageHMACString)
	default:
		return "", "", nil, errors.New("'" + HMACFormat + "' not a supported format")
	}

	if err != nil {
		return "", "", nil, err
	}

	return message, time, messageHMAC, nil
}

func AuthRequest(r *http.Request, userID uint) (bool, error) {
	message, timeString, messageHMAC, err := parseRequest(r)
	if err != nil {
		return false, err
	}

	timeInt, err := strconv.Atoi(timeString)
	if err != nil {
		return false, errors.New("Error parsing time (seconds since epoch)")
	}

	nowMillis := time.Now().Unix()
	delay := int64(timeInt) - nowMillis

	// rejects if times are more than 10 seconds apart
	if delay < -60 || delay > 60 {
		return false, errors.New("Time difference too large")
	}

	secret, ok := secretMap[userID]
	if !ok {
		return false, errors.New("No secret for that user")
	}

	secretString := secret.String()
	return checkMAC(message, secretString, messageHMAC), nil
}
