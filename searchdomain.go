// odoojrpc - go library to access Odoo server via Json RPC
// Copyright (C) 2021  Peter Preeper

// This library is free software; you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 2.1 of the License, or (at your option) any later version.

// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public
// License along with this library; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301
// USA
package odoojrpc

import (
	"errors"
	"regexp"
	"sort"
	"strings"
)

type StringNode struct {
	Start int
	End   int
}

var (
	errSyntax     = errors.New("invalid syntax")
	comparators   = `(=|!=|>|>=|<|<=|=?|=like|like|not like|ilike|not ilike|=ilike|in|not in|child_of|parent_of)`
	baseArg       = `(\(\s*'\w+'\s*,\s*'` + comparators + `'\s*,\s*'(\w|\s)+'\s*\))`
	sTerm         = `^` + baseArg + `$`
	andorTerms    = `(('&'|'\|')\s*,\s*` + baseArg + `\s*,\s*` + baseArg + `)`
	notTerms      = `('!'\s*,\s*` + baseArg + `)`
	combinedTerms = `(,\s*(` + baseArg + `|` + andorTerms + `|` + notTerms + `))*`
	mTerm         = `^(\[\s*` + baseArg + combinedTerms + `\s*\])$`
	argCheck      = mTerm + "|" + sTerm
)

var (
	ac      = regexp.MustCompile(argCheck)
	reSTerm = regexp.MustCompile(sTerm)
	reMTerm = regexp.MustCompile(mTerm)
)

var (
	reBase  = regexp.MustCompile(baseArg)
	reAddOr = regexp.MustCompile(andorTerms)
	reNot   = regexp.MustCompile(notTerms)
)

var (
	reAndOrTerm = regexp.MustCompile(`^\s*('&'|'\|')`)
	reNotTerm   = regexp.MustCompile(`^\s*('!')`)
)

// term searches
var (
	fieldTerm      = regexp.MustCompile(`\(\s*'(\w)+'`)
	comparatorTerm = regexp.MustCompile(`\s*,'` + comparators + `'\s*,`)
	valueTerm      = regexp.MustCompile(`\s*'(\w|\s)+'\s*\)`)
)

func SearchDomain(domain string) (filter FilterArg, err error) {
	if domain == "" {
		return FilterArg{}, nil
	}

	if !ac.MatchString(domain) {
		return FilterArg{}, errSyntax
	}

	// single term
	if reSTerm.MatchString(domain) {
		return FilterArg{patternSplit(domain)}, nil
	}

	// multi term
	if reMTerm.MatchString(domain) {
		aa := []StringNode{}
		nn := []StringNode{}
		bb := []StringNode{}

		for _, v := range reAddOr.FindAllStringIndex(domain, -1) {
			aa = append(aa, StringNode{Start: v[0], End: v[1]})
		}

		for _, v := range reNot.FindAllStringIndex(domain, -1) {
			nn = append(nn, StringNode{Start: v[0], End: v[1]})
		}

		for _, v := range reBase.FindAllStringIndex(domain, -1) {
			bb = append(bb, StringNode{Start: v[0], End: v[1]})
		}

		nl := buildNodeList(aa, nn, bb)
		for _, n := range nl {
			ss := stringSplit(domain, n.Start, n.End)
			filter = append(filter, patternSplit(ss))
		}
		return filter, nil
	}
	return filter, nil
}

func buildNodeList(aa, nn, bb []StringNode) (nl []StringNode) {
	nl = insideNodeList(nl, aa)
	nl = insideNodeList(nl, nn)
	nl = insideNodeList(nl, bb)
	sort.Slice(nl, func(i int, j int) bool {
		return nl[i].Start < nl[j].Start
	})
	return
}

func insideNodeList(aa, bb []StringNode) []StringNode {
	for _, n := range bb {
		inside := false
		for _, a := range aa {
			if n.Start >= a.Start && n.End <= a.End {
				inside = true
			}
		}
		if !inside {
			aa = append(aa, n)
		}
	}
	return aa
}

func stringSplit(term string, start, end int) string {
	b := []byte(term)
	bStr := b[start:end]
	return string(bStr)
}

func patternSplit(statement string) FilterArg {
	if reAddOr.MatchString(statement) {
		opCondition := reAndOrTerm.FindAllString(statement, -1)
		op := termTrimQuote(opCondition[0])
		terms := reBase.FindAllString(statement, -1)
		t1 := termSplit(terms[0])
		t2 := termSplit(terms[1])
		return FilterArg{op, FilterArg{t1[0], t1[1], t1[2]}, FilterArg{t2[0], t2[1], t2[2]}}
	}

	if reNot.MatchString(statement) {
		opCondition := reNotTerm.FindAllString(statement, -1)
		op := termTrimQuote(opCondition[0])
		terms := reBase.FindAllString(statement, -1)
		tt := termSplit(terms[0])
		return FilterArg{op, FilterArg{tt[0], tt[1], tt[2]}}
	}

	if reBase.MatchString(statement) {
		terms := reBase.FindAllString(statement, -1)
		tt := termSplit(terms[0])
		return FilterArg{tt[0], tt[1], tt[2]}
	}
	return FilterArg{}
}

func termSplit(term string) []string {
	field := fieldTerm.FindAllString(term, -1)
	fieldStr := termTrimQuote(strings.Trim(field[0], "("))

	comparator := comparatorTerm.FindAllString(term, -1)
	comparatorStr := termTrimQuote(strings.Trim(comparator[0], ","))

	value := valueTerm.FindAllString(term, -1)
	valueStr := termTrimQuote(strings.Trim(value[0], ")"))

	return []string{fieldStr, comparatorStr, valueStr}
}

func termTrimQuote(term string) string {
	return strings.TrimSpace(strings.Trim(term, "'"))
}
