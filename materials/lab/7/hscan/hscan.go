package hscan

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
)

//==========================================================================\\

var shalookup = make(map[string]string)
var md5lookup = make(map[string]string)
var mutexSHA sync.Mutex
var mutexMD5 sync.Mutex

func GuessSingle(sourceHash string, filename string) string {

	hashLen := len(sourceHash)

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		password := scanner.Text()

		// TODO - From the length of the hash you should know which one of these to check ...
		// add a check and logicial structure

		if hashLen == 32 {
			hash := fmt.Sprintf("%x", md5.Sum([]byte(password)))
			if hash == sourceHash {
				fmt.Printf("[+] Password found (MD5): %s\n", password)
				return password
			}
		} else if hashLen == 64 {
			hash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
			if hash == sourceHash {
				fmt.Printf("[+] Password found (SHA-256): %s\n", password)
				return password
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
	return ""
}


func GenHashMaps(filename string) {
	//TODO
	//itterate through a file (look in the guessSingle function above)
	//rather than check for equality add each hash:passwd entry to a map SHA and MD5 where the key = hash and the value = password
	//TODO at the very least use go subroutines to generate the sha and md5 hashes at the same time
	//OPTIONAL -- Can you use workers to make this even faster

	//TODO create a test in hscan_test.go so that you can time the performance of your implementation
	//Test and record the time it takes to scan to generate these Maps
	// 1. With and without using go subroutines
	// 2. Compute the time per password (hint the number of passwords for each file is listed on the site...)

	var routines sync.WaitGroup
	
	f, err := os.Open(filename)
	g, err1 := os.Open(filename)
	if err != nil || err1 != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	defer g.Close()

	routines.Add(1)
	go StoreSHAHash(f, &routines)
	routines.Add(1)
	go StoreMD5Hash(g, &routines)
	routines.Wait()
}

func StoreMD5Hash(file *os.File, routines *sync.WaitGroup){
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		password := scanner.Text()
		mutexMD5.Lock()
		hash := fmt.Sprintf("%x", md5.Sum([]byte(password)))
		md5lookup[hash] = password
		mutexMD5.Unlock()
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
	routines.Done()
}

func StoreSHAHash(file *os.File, routines *sync.WaitGroup){
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		password := scanner.Text()
		mutexSHA.Lock()
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
		shalookup[hash] = password
		mutexSHA.Unlock()
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
	routines.Done()
}

func GetSHA(hash string) (string, error) {
	password, ok := shalookup[hash]
	if ok {
		return password, nil

	} else {

		return "", errors.New("password does not exist")

	}
}

//TODO
func GetMD5(hash string) (string, error) {
	password, ok := md5lookup[hash]
	if ok {
		return password, nil

	} else {

		return "", errors.New("password does not exist")

	}
}
