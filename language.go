package main

import (
	"io"
	"math/rand"
	"os"
	"strconv"
	"sync"
)

var languageError error
var languageOnce sync.Once

type LanguageWord struct {
	ID        string
	Noun      *LanguageNoun
	Prefix    *LanguagePrefix
	Verb      *LanguageVerb
	Adjective *LanguageAdjective

	Translation struct {
		Dwarf  string
		Human  string
		Goblin string
		Elf    string
	}
}

type LanguageNoun struct {
	Singular string
	Plural   string

	TheSingular           bool
	ThePlural             bool
	TheCompoundSingular   bool
	TheCompoundPlural     bool
	OfSingular            bool
	OfPlural              bool
	FrontCompoundSingular bool
	FrontCompoundPlural   bool
	RearCompoundSingular  bool
	RearCompoundPlural    bool
}

type LanguagePrefix struct {
	Prefix string

	FrontCompound bool
	TheCompound   bool
}

type LanguageVerb struct {
	PresentFirst string
	PresentThird string
	Preterite    string
	PastPart     string
	PresentPart  string

	Standard bool
}

type LanguageAdjective struct {
	Adjective string
	Distance  int

	TheCompound   bool
	FrontCompound bool
	RearCompound  bool
}

var languageWords []*LanguageWord
var languageWordsMap = make(map[string]*LanguageWord)

var languageNouns []*LanguageWord
var languageFront []*LanguageWord
var languageRear []*LanguageWord

func languageInit() {
	f, err := os.Open("raws/objects/language_words.txt")
	if err != nil {
		languageError = err
		return
	}
	defer f.Close()

	t := NewRawsTokenizer(f)
	for {
		tok, err := t.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			languageError = err
			return
		}
		switch tok[0] {
		case "OBJECT":
		case "WORD":
			languageWords = append(languageWords, &LanguageWord{
				ID: tok[1],
			})
			languageWordsMap[tok[1]] = languageWords[len(languageWords)-1]
		case "NOUN":
			languageWords[len(languageWords)-1].Noun = &LanguageNoun{
				Singular: tok[1],
				Plural:   tok[2],
			}
		case "PREFIX":
			languageWords[len(languageWords)-1].Prefix = &LanguagePrefix{
				Prefix: tok[1],
			}
		case "VERB":
			languageWords[len(languageWords)-1].Verb = &LanguageVerb{
				PresentFirst: tok[1],
				PresentThird: tok[2],
				Preterite:    tok[3],
				PastPart:     tok[4],
				PresentPart:  tok[5],
			}
		case "ADJ":
			languageWords[len(languageWords)-1].Adjective = &LanguageAdjective{
				Adjective: tok[1],
			}
		case "ADJ_DIST":
			languageWords[len(languageWords)-1].Adjective.Distance, err = strconv.Atoi(tok[1])
			if err != nil {
				languageError = err
				return
			}
		case "FRONT_COMPOUND_NOUN_SING":
			languageWords[len(languageWords)-1].Noun.FrontCompoundSingular = true
		case "FRONT_COMPOUND_NOUN_PLUR":
			languageWords[len(languageWords)-1].Noun.FrontCompoundPlural = true
		case "REAR_COMPOUND_NOUN_SING":
			languageWords[len(languageWords)-1].Noun.RearCompoundSingular = true
		case "REAR_COMPOUND_NOUN_PLUR":
			languageWords[len(languageWords)-1].Noun.RearCompoundPlural = true
		case "THE_COMPOUND_NOUN_SING":
			languageWords[len(languageWords)-1].Noun.TheCompoundSingular = true
		case "THE_COMPOUND_NOUN_PLUR":
			languageWords[len(languageWords)-1].Noun.TheCompoundPlural = true
		case "THE_NOUN_SING":
			languageWords[len(languageWords)-1].Noun.TheSingular = true
		case "THE_NOUN_PLUR":
			languageWords[len(languageWords)-1].Noun.ThePlural = true
		case "OF_NOUN_SING":
			languageWords[len(languageWords)-1].Noun.OfSingular = true
		case "OF_NOUN_PLUR":
			languageWords[len(languageWords)-1].Noun.OfPlural = true
		case "FRONT_COMPOUND_PREFIX":
			languageWords[len(languageWords)-1].Prefix.FrontCompound = true
		case "THE_COMPOUND_PREFIX":
			languageWords[len(languageWords)-1].Prefix.TheCompound = true
		case "STANDARD_VERB":
			languageWords[len(languageWords)-1].Verb.Standard = true
		case "FRONT_COMPOUND_ADJ":
			languageWords[len(languageWords)-1].Adjective.FrontCompound = true
		case "REAR_COMPOUND_ADJ":
			languageWords[len(languageWords)-1].Adjective.RearCompound = true
		case "THE_COMPOUND_ADJ":
			languageWords[len(languageWords)-1].Adjective.TheCompound = true
		default:
			panic(tok[0])
		}
	}

	for _, l := range []struct {
		L string
		F func(*LanguageWord, string)
	}{
		{"DWARF", func(w *LanguageWord, s string) { w.Translation.Dwarf = s }},
		{"HUMAN", func(w *LanguageWord, s string) { w.Translation.Human = s }},
		{"GOBLIN", func(w *LanguageWord, s string) { w.Translation.Goblin = s }},
		{"ELF", func(w *LanguageWord, s string) { w.Translation.Elf = s }},
	} {
		f, err := os.Open("raws/objects/language_" + l.L + ".txt")
		if err != nil {
			languageError = err
			return
		}
		defer f.Close()

		t := NewRawsTokenizer(f)
		for {
			tok, err := t.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				languageError = err
				return
			}
			switch tok[0] {
			case "OBJECT":
			case "TRANSLATION":
			case "T_WORD":
				l.F(languageWordsMap[tok[1]], tok[2])
			default:
				panic(tok[0])
			}
		}
	}

	for _, w := range languageWords {
		if w.Noun != nil {
			languageNouns = append(languageNouns, w)
		}

		if (w.Noun != nil && (w.Noun.FrontCompoundSingular || w.Noun.FrontCompoundPlural)) || (w.Prefix != nil && w.Prefix.FrontCompound) || (w.Verb != nil && w.Verb.Standard) || (w.Adjective != nil && w.Adjective.FrontCompound) {
			languageFront = append(languageFront, w)
		}

		if (w.Noun != nil && (w.Noun.RearCompoundSingular || w.Noun.RearCompoundPlural)) || (w.Verb != nil && w.Verb.Standard) {
			languageRear = append(languageRear, w)
		}
	}
}

func LanguageInit() error {
	languageOnce.Do(languageInit)
	return languageError
}

func GenerateNameParts(r *rand.Rand) (first, front, rear *LanguageWord) {
	first = languageNouns[r.Intn(len(languageNouns))]
	front = languageFront[r.Intn(len(languageFront))]
	rear = languageRear[r.Intn(len(languageRear))]
	return
}
