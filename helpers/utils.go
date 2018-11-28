package helpers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// CheckEnv Check the Azure creds are set in the environment variables
func CheckEnv() (err error) {

	azureCreds := []string{
		"AZURE_TENANT_ID",
		"AZURE_CLIENT_ID",
		"AZURE_CLIENT_SECRET",
		"AZURE_SUBSCRIPTION_ID",
	}

	for _, cred := range azureCreds {
		if os.Getenv(cred) == "" {
			log.Printf("credential variable " + cred + " has not be set")
			err = errors.New("error, missing envrionment variables. run `az ad sp create-for-rbac -n \"<yourAppName>\"' -o json --sdk-auth to create a service principal and generate the necessary credential variables")
		} else {
			log.Printf("%v variable was found. OK.", cred)
		}
	}

	return err
}

// FatalError Quits if error is fatal
func FatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//PrintError Prints if error
func PrintError(err error) {
	if err != nil {
		log.Printf(err.Error())
	}
}

// ReadJSON Reads json and returns a map
func ReadJSON(path string) (*map[string]interface{}, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	contents := make(map[string]interface{})
	json.Unmarshal(data, &contents)
	return &contents, nil
}

// RandStringBytesMaskImprSrc Generates a random string that is  'n' strings longs
func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
