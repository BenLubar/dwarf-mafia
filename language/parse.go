package language

import (
	"io"
	"math/rand"
	"os"
	"strconv"

	"github.com/BenLubar/df2014/raws"
)

type Word struct {
	ID        string
	Noun      *Noun
	Prefix    *Prefix
	Verb      *Verb
	Adjective *Adjective

	Translation struct {
		Dwarf  string
		Human  string
		Goblin string
		Elf    string
	}
}

type Noun struct {
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

type Prefix struct {
	Prefix string

	FrontCompound bool
	TheCompound   bool
}

type Verb struct {
	PresentFirst string
	PresentThird string
	Preterite    string
	PastPart     string
	PresentPart  string

	Standard bool
}

type Adjective struct {
	Adjective string
	Distance  int

	TheCompound   bool
	FrontCompound bool
	RearCompound  bool
}

var words []*Word
var wordsMap = make(map[string]*Word)

var nouns, fronts, rears []*Word

func init() {
	f, err := os.Open("raws/objects/language_words.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	t := raws.NewTokenizer(f)
	for {
		tok, err := t.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		switch tok[0] {
		case "OBJECT":
		case "WORD":
			w := &Word{
				ID: tok[1],
			}
			words = append(words, w)
			wordsMap[tok[1]] = w
		case "NOUN":
			w := words[len(words)-1]
			w.Noun = &Noun{
				Singular: tok[1],
				Plural:   tok[2],
			}
		case "PREFIX":
			w := words[len(words)-1]
			w.Prefix = &Prefix{
				Prefix: tok[1],
			}
		case "VERB":
			w := words[len(words)-1]
			w.Verb = &Verb{
				PresentFirst: tok[1],
				PresentThird: tok[2],
				Preterite:    tok[3],
				PastPart:     tok[4],
				PresentPart:  tok[5],
			}
		case "ADJ":
			w := words[len(words)-1]
			w.Adjective = &Adjective{
				Adjective: tok[1],
			}
		case "ADJ_DIST":
			w := words[len(words)-1]
			w.Adjective.Distance, err = strconv.Atoi(tok[1])
			if err != nil {
				panic(err)
			}
		case "FRONT_COMPOUND_NOUN_SING":
			words[len(words)-1].Noun.FrontCompoundSingular = true
		case "FRONT_COMPOUND_NOUN_PLUR":
			words[len(words)-1].Noun.FrontCompoundPlural = true
		case "REAR_COMPOUND_NOUN_SING":
			words[len(words)-1].Noun.RearCompoundSingular = true
		case "REAR_COMPOUND_NOUN_PLUR":
			words[len(words)-1].Noun.RearCompoundPlural = true
		case "THE_COMPOUND_NOUN_SING":
			words[len(words)-1].Noun.TheCompoundSingular = true
		case "THE_COMPOUND_NOUN_PLUR":
			words[len(words)-1].Noun.TheCompoundPlural = true
		case "THE_NOUN_SING":
			words[len(words)-1].Noun.TheSingular = true
		case "THE_NOUN_PLUR":
			words[len(words)-1].Noun.ThePlural = true
		case "OF_NOUN_SING":
			words[len(words)-1].Noun.OfSingular = true
		case "OF_NOUN_PLUR":
			words[len(words)-1].Noun.OfPlural = true
		case "FRONT_COMPOUND_PREFIX":
			words[len(words)-1].Prefix.FrontCompound = true
		case "THE_COMPOUND_PREFIX":
			words[len(words)-1].Prefix.TheCompound = true
		case "STANDARD_VERB":
			words[len(words)-1].Verb.Standard = true
		case "FRONT_COMPOUND_ADJ":
			words[len(words)-1].Adjective.FrontCompound = true
		case "REAR_COMPOUND_ADJ":
			words[len(words)-1].Adjective.RearCompound = true
		case "THE_COMPOUND_ADJ":
			words[len(words)-1].Adjective.TheCompound = true
		default:
			panic(tok[0])
		}
	}

	for _, l := range []struct {
		L string
		F func(*Word, string)
	}{
		{"DWARF", func(w *Word, s string) { w.Translation.Dwarf = s }},
		{"HUMAN", func(w *Word, s string) { w.Translation.Human = s }},
		{"GOBLIN", func(w *Word, s string) { w.Translation.Goblin = s }},
		{"ELF", func(w *Word, s string) { w.Translation.Elf = s }},
	} {
		func() {
			f, err := os.Open("raws/objects/language_" + l.L + ".txt")
			if err != nil {
				panic(err)
			}
			defer f.Close()

			t := raws.NewTokenizer(f)
			for {
				tok, err := t.Next()
				if err == io.EOF {
					break
				}
				if err != nil {
					panic(err)
				}
				switch tok[0] {
				case "OBJECT":
				case "TRANSLATION":
				case "T_WORD":
					l.F(wordsMap[tok[1]], tok[2])
				default:
					panic(tok[0])
				}
			}
		}()
	}

	for _, w := range words {
		if w.Noun != nil {
			nouns = append(nouns, w)
		}

		if (w.Noun != nil && (w.Noun.FrontCompoundSingular || w.Noun.FrontCompoundPlural)) || (w.Prefix != nil && w.Prefix.FrontCompound) || (w.Verb != nil && w.Verb.Standard) || (w.Adjective != nil && w.Adjective.FrontCompound) {
			fronts = append(fronts, w)
		}

		if (w.Noun != nil && (w.Noun.RearCompoundSingular || w.Noun.RearCompoundPlural)) || (w.Verb != nil && w.Verb.Standard) {
			rears = append(rears, w)
		}
	}
}

func GenerateNameParts(r *rand.Rand) (first, front, rear *Word) {
	first = nouns[r.Intn(len(nouns))]
	front = fronts[r.Intn(len(fronts))]
	rear = rears[r.Intn(len(rears))]
	return
}
