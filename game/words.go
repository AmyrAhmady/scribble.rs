package game

import (
	"math/rand"
	"strings"
	"time"

	"github.com/gobuffalo/packr/v2"
	"golang.org/x/text/cases"
)

var (
	wordListCache       = make(map[string][]string)
	languageIdentifiers = map[string]string{
		"english": "en",
		"italian": "it",
		"german":  "de",
		"french":  "fr",
		"dutch":   "nl",
	}
	wordBox = packr.New("words", "../resources/words")
)

func getLanguageIdentifier(language string) string {
	return languageIdentifiers[language]
}

func readWordList(lowercaser cases.Caser, chosenLanguage string) ([]string, error) {
	languageIdentifier := getLanguageIdentifier(chosenLanguage)
	list, available := wordListCache[languageIdentifier]
	if available {
		return list, nil
	}

	wordListFile, pkgerError := wordBox.FindString(languageIdentifier)
	if pkgerError != nil {
		panic(pkgerError)
	}

	tempWords := strings.Split(wordListFile, "\n")
	var words []string
	for _, word := range tempWords {
		word = strings.TrimSpace(word)

		//Newlines will just be empty strings
		if word == "" {
			continue
		}

		//The "i" was "impossible", as in "impossible to draw", tag initially supplied.
		if strings.HasSuffix(word, "#i") {
			continue
		}

		//Since not all words use the tag system, we can just instantly return for words that don't use it.
		lastIndexNumberSign := strings.LastIndex(word, "#")
		if lastIndexNumberSign == -1 {
			words = append(words, lowercaser.String(word))
		} else {
			words = append(words, lowercaser.String(word[:lastIndexNumberSign]))
		}
	}

	wordListCache[languageIdentifier] = words

	return words, nil
}

// GetRandomWords gets 3 random words for the passed Lobby. The words will be
// chosen from the custom words and the default dictionary, depending on the
// settings specified by the Lobby-owner.
func GetRandomWords(lobby *Lobby) []string {
	rand.Seed(time.Now().Unix())
	wordsNotToPick := lobby.alreadyUsedWords
	wordOne := getRandomWordWithCustomWordChance(lobby, wordsNotToPick, lobby.CustomWords, lobby.CustomWordsChance)
	wordsNotToPick = append(wordsNotToPick, wordOne)
	wordTwo := getRandomWordWithCustomWordChance(lobby, wordsNotToPick, lobby.CustomWords, lobby.CustomWordsChance)
	wordsNotToPick = append(wordsNotToPick, wordTwo)
	wordThree := getRandomWordWithCustomWordChance(lobby, wordsNotToPick, lobby.CustomWords, lobby.CustomWordsChance)

	return []string{
		wordOne,
		wordTwo,
		wordThree,
	}
}

func getRandomWordWithCustomWordChance(lobby *Lobby, wordsAlreadyUsed []string, customWords []string, customWordChance int) string {
	if len(lobby.CustomWords) > 0 && customWordChance > 0 && rand.Intn(100)+1 <= customWordChance {
		return getUnusedCustomWord(lobby, wordsAlreadyUsed, customWords)
	}

	return getUnusedRandomWord(lobby, wordsAlreadyUsed)
}

func getUnusedCustomWord(lobby *Lobby, wordsAlreadyUsed []string, customWords []string) string {
OUTER_LOOP:
	for _, word := range customWords {
		for _, usedWord := range wordsAlreadyUsed {
			if usedWord == word {
				continue OUTER_LOOP
			}
		}

		return word
	}

	return getUnusedRandomWord(lobby, wordsAlreadyUsed)
}

func getUnusedRandomWord(lobby *Lobby, wordsAlreadyUsed []string) string {
	//We attempt to find a random word for a hundred times, afterwards we just use any.
	randomnessAttempts := 0
	var word string
OUTER_LOOP:
	for {
		word = lobby.words[rand.Int()%len(lobby.words)]
		for _, usedWord := range wordsAlreadyUsed {
			if usedWord == word {
				if randomnessAttempts == 100 {
					break OUTER_LOOP
				}

				randomnessAttempts++
				continue OUTER_LOOP
			}
		}
		break
	}

	return word
}
