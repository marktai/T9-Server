package auth_test

import (
	"auth"
)

// func MakeUser(user, pass string) (uint, error) {
// 	_, err := getUserID(user)
// 	if err == nil {
// 		return 0, errors.New("User already made")
// 	}

// 	id, err := getUniqueID()
// 	if err != nil {
// 		return 0, err
// 	}

// 	saltHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
// 	if err != nil {
// 		return 0, err
// 	}
// 	saltHashString := base64.StdEncoding.EncodeToString(saltHash)

// 	err = db.Db.Ping()
// 	if err != nil {
// 		return 0, err
// 	}

// 	addUser, err := db.Db.Prepare("INSERT INTO users VALUES(?, ?, ?)")
// 	if err != nil {
// 		return 0, err
// 	}

// 	_, err = addUser.Exec(id, user, saltHashString)

// 	if err != nil {
// 		return 0, err
// 	}

// 	return id, nil

// }

// func getUserID(user string) (uint, error) {
// 	var userID uint
// 	err := db.Db.QueryRow("SELECT id FROM users WHERE name=?", user).Scan(&userID)

// 	return userID, err
// }

// func getSaltHash(userID uint) ([]byte, error) {
// 	saltHashString := ""
// 	err := db.Db.QueryRow("SELECT salthash FROM users WHERE id=?", userID).Scan(&saltHashString)
// 	if err != nil {
// 		return nil, err
// 	}
// 	saltHash, err := base64.StdEncoding.DecodeString(saltHashString)
// 	return saltHash, err
// }

// func Login(user, pass string) (uint, *Secret, error) {

// 	userID, err := getUserID(user)
// 	if err != nil {
// 		return 0, nil, err
// 	}

// 	hash, err := getSaltHash(userID)
// 	if err != nil {
// 		return 0, nil, err
// 	}

// 	err = bcrypt.CompareHashAndPassword(hash, []byte(pass))

// 	if err != nil {
// 		return 0, nil, err
// 	}

// 	if _, ok := secretMap[userID]; !ok || secretMap[userID].Expired() {
// 		secret, err := newSecret()
// 		if err != nil {
// 			return 0, nil, err
// 		}
// 		secretMap[userID] = secret
// 	}

// 	secretMap[userID].resetExpiration()

// 	return userID, secretMap[userID], nil
// }

// func VerifySecret(user, inpSecret string) (uint, *Secret, error) {

// 	userID, err := getUserID(user)
// 	if err != nil {
// 		return 0, nil, err
// 	}

// 	if _, ok := secretMap[userID]; !ok {
// 		return 0, nil, errors.New("No secret found for user")
// 	} else if secretMap[userID].Expired() {
// 		return 0, nil, errors.New("Secret has expired")
// 	} else if secretMap[userID].String() != inpSecret {
// 		return 0, nil, errors.New("Secrets do not match")
// 	}

// 	secretMap[userID].resetExpiration()

// 	return userID, secretMap[userID], nil
// }

// func ComputeHmac256(message, key string) []byte {
// 	mac := hmac.New(sha256.New, []byte(key))
// 	mac.Write([]byte(message))
// 	expectedMAC := mac.Sum(nil)

// 	return expectedMAC
// }

// // CheckMAC reports whether messageHMAC is a valid HMAC tag for message.
// func checkMAC(key, message string, messageHMAC []byte) bool {
// 	expectedMAC := ComputeHmac256(message, key)
// 	return hmac.Equal(messageHMAC, expectedMAC)
// }

// // returns userID, message used to generate HMAC, and HMAC from request
// func parseRequest(r *http.Request) (string, string, []byte, error) {
// 	timeSlice, ok := r.Header["Time-Sent"]
// 	if !ok || timeSlice == nil || len(timeSlice) == 0 {
// 		return "", "", nil, errors.New("No Time-Sent header provided")
// 	}

// 	time := timeSlice[0]

// 	message := time + ":" + r.URL.String()

// 	messageHMACSlice, ok := r.Header["Hmac"]
// 	if !ok || messageHMACSlice == nil || len(messageHMACSlice) == 0 {
// 		return "", "", nil, errors.New("No HMAC header provided")
// 	}

// 	messageHMACString := messageHMACSlice[0]
// 	HMACEncoding := ""

// 	encodingSlice, ok := r.Header["Encoding"]
// 	if ok && encodingSlice != nil {
// 		encoding := encodingSlice[0]
// 		if ok {
// 			switch strings.ToLower(encoding) {
// 			case "base64", "64":
// 				HMACEncoding = "base64"
// 			case "hex", "hexadecimal":
// 				HMACEncoding = "hex"
// 			case "binary", "bits":
// 				HMACEncoding = "binary"
// 			case "decimal":
// 				HMACEncoding = "decimal"
// 			default:
// 				HMACEncoding = encoding
// 			}
// 		}
// 	} else {
// 		HMACEncoding = "hex"
// 	}

// 	var messageHMAC []byte
// 	var err error
// 	switch HMACEncoding {
// 	case "base64":
// 		messageHMAC, err = base64.StdEncoding.DecodeString(messageHMACString)
// 	case "hex":
// 		messageHMAC, err = hex.DecodeString(messageHMACString)
// 	default:
// 		return "", "", nil, errors.New("'" + HMACEncoding + "' not a supported encoding")
// 	}

// 	if err != nil {
// 		return "", "", nil, err
// 	}

// 	return message, time, messageHMAC, nil
// }

// // Check README.md for documentation
// // Request Headers
// // HMAC - encoded HMAC with SHA 256
// // Encoding - encoding format (default hex)
// // Time-Sent - seconds since epoch

// // Verifies whether a request is correctly authorized
// func AuthRequest(r *http.Request, userID uint) (bool, error) {
// 	message, timeString, messageHMAC, err := parseRequest(r)
// 	if err != nil {
// 		return false, err
// 	}

// 	timeInt, err := strconv.Atoi(timeString)
// 	if err != nil {
// 		return false, errors.New("Error parsing time (seconds since epoch)")
// 	}

// 	delay := int64(timeInt) - time.Now().Unix()

// 	// rejects if times are more than 10 seconds apart
// 	if delay < -10 || delay > 10 {
// 		return false, errors.New("Time difference too large")
// 	}

// 	secret, ok := secretMap[userID]
// 	if !ok {
// 		return false, errors.New("No secret for that user")
// 	}

// 	if secret.Expired() {
// 		return false, errors.New("Secret expired")
// 	}

// 	secretString := secret.String()
// 	return checkMAC(secretString, message, messageHMAC), nil
// }
