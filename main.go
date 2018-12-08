package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

func main() {
	demo := flag.Bool("demo", false, "")
	in := flag.String("in", "./in", "")
	out := flag.String("out", "./out/results.zip", "")
	flag.Parse()
	participants := readNames(*in)
	pairs := findPairs(randomize(participants))
	if *demo {
		print(pairs)
		return
	}
	os.MkdirAll(filepath.Dir(*out), os.ModePerm)
	writeResults(*in, *out, pairs)
}

func print(pairs map[string]string) {
	for k, v := range pairs {
		fmt.Printf("%s -> %s\n", k, v)
	}
}

func findPairs(participants []string) map[string]string {
	pairs := map[string]string{}
	takers := map[string]bool{}
	for _, giver := range participants {
	inner:
		for _, taker := range participants {
			_, alreadyGave := pairs[giver]
			givingToTaker := pairs[taker] == giver
			alreadyTook := takers[taker]
			if !alreadyGave && !alreadyTook && !givingToTaker && giver != taker {
				takers[taker] = true
				pairs[giver] = taker
				break inner
			}
		}
	}
	return pairs
}

func randomize(participants []string) []string {
	s := rand.NewSource(time.Now().UTC().UnixNano())
	r := rand.New(s)
	count := len(participants)
	permutation := r.Perm(count)
	randomized := make([]string, count)
	for i, p := range permutation {
		randomized[i] = participants[p]
	}
	return randomized
}

func readNames(dir string) []string {
	var names = make([]string, 0)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		names = append(names, file.Name())
	}
	return names
}

func encrypt(plain []byte, keyPath string) []byte {
	sshKeygen := exec.Command("ssh-keygen", "-f", keyPath, "-e", "-m", "PKCS8")
	var stdout bytes.Buffer
	sshKeygen.Stdout = &stdout
	err := sshKeygen.Run()
	if err != nil {
		log.Fatalf("Error running ssh-keygen: %v\n", err)
	}
	pem := stdout.Bytes()
	stdout.Reset()
	tempFilename := uuid.New().String()
	err = ioutil.WriteFile(tempFilename, pem, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tempFilename)
	openssl := exec.Command("openssl", "rsautl", "-encrypt", "-pubin", "-inkey", tempFilename, "-ssl")
	openssl.Stdin = bytes.NewBuffer(plain)
	openssl.Stdout = &stdout
	err = openssl.Run()
	if err != nil {
		log.Fatal(err)
	}
	return stdout.Bytes()
}

var readme = "# How to Read the Result\n\n" +
	"In the directory containing this README:\n\n" +
	"```bash\n" +
	"chmod +x decrypt\n" +
	"cat <YOUR NAME> | ./decrypt\n" +
	"```\n\n" +
	"By default the private key located in ~/.ssh/id_rsa is used.\n"

var decryptScript = "#!/bin/bash\n\n" +
	"cat | openssl rsautl -decrypt -inkey ~/.ssh/id_rsa\n"

func writeToArchive(w *zip.Writer, name string, content []byte) {
	f, err := w.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Write(content)
	if err != nil {
		log.Fatal(err)
	}
}

func writeResults(sourceDir string, targetPath string, pairs map[string]string) {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	writeToArchive(w, "README.md", []byte(readme))
	writeToArchive(w, "decrypt", []byte(decryptScript))
	for k, v := range pairs {
		sourcePath := fmt.Sprintf("%s/%s", sourceDir, k)
		encrypted := encrypt([]byte(v), sourcePath)
		f, err := w.Create(k)
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.Write(encrypted)
		if err != nil {
			log.Fatal(err)
		}
	}
	err := w.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(targetPath, buf.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
