package main

import (
	"math/rand"
	"strings"

	"github.com/BenLubar/dwarf-mafia/language"
)

func GenerateName(r *rand.Rand) (english, dwarf string) {
	first, front, rear := language.GenerateNameParts(r)

	if r.Intn(3) <= 1 && first.Noun.Singular != "" {
		english = first.Noun.Singular + " "
	} else if first.Noun.Plural != "" {
		english = first.Noun.Plural + " "
	} else {
		english = first.Noun.Singular + " "
	}
	dwarf = first.Translation.Dwarf + " "

	var possible []string

	if front.Noun != nil && front.Noun.FrontCompoundSingular {
		possible = append(possible, front.Noun.Singular)
	}
	if front.Noun != nil && front.Noun.FrontCompoundPlural {
		possible = append(possible, front.Noun.Plural)
	}
	if front.Prefix != nil && front.Prefix.FrontCompound {
		possible = append(possible, front.Prefix.Prefix)
	}
	if front.Verb != nil && front.Verb.Standard {
		possible = append(possible, front.Verb.PresentFirst, front.Verb.PresentThird, front.Verb.Preterite)
	}
	if front.Adjective != nil && front.Adjective.FrontCompound {
		possible = append(possible, front.Adjective.Adjective)
	}

	english += possible[r.Intn(len(possible))]
	dwarf += front.Translation.Dwarf

	possible = possible[:0]

	if rear.Noun != nil && rear.Noun.RearCompoundSingular {
		possible = append(possible, rear.Noun.Singular)
	}
	if rear.Noun != nil && rear.Noun.RearCompoundPlural {
		possible = append(possible, rear.Noun.Plural)
	}
	if rear.Verb != nil && rear.Verb.Standard {
		possible = append(possible, rear.Verb.PresentFirst, rear.Verb.PresentThird, rear.Verb.Preterite)
	}

	english += possible[r.Intn(len(possible))]
	dwarf += rear.Translation.Dwarf

	english = strings.Title(english)
	dwarf = strings.Title(dwarf)

	return
}
