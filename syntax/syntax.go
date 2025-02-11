// Package syntax defines the syntax for specifying packages and platforms.
//
// A package can either be specified by a string in the format of
// "platform/name@version". Only the name is required, both platform and version
// can be omitted.
//
// Valid Examples:
//   - carpet
//   - mcdr/prime-backup
//   - fabric/jade@1.0.0
//   - fabric@12.0
//   - minecraft@1.19 (recommended)
//   - minecraft/minecraft@1.16.5 (= minecraft@1.16.5)
//   - 1.8.9 (= minecraft@1.8.9)
package syntax

import (
	"errors"
	"log"
	"lucy/syntaxtypes"
	"strings"
)

// sanitize tolerates some common interchangeability between characters. This
// includes underscores, chinese full stops, and backslashes. It also converts
// uppercase characters to lowercase.
func sanitize(s string) (clean string) {
	clean = ""
	for _, char := range s {
		if char == '_' {
			clean += string('-')
		} else if char == '\\' {
			clean += string('/')
		} else if char == '。' {
			clean += string('.')
		} else if char >= 'A' && char <= 'Z' {
			clean += strings.ToLower(string(char))
		} else {
			clean += string(char)
		}
	}
	return
}

var (
	ESyntax   = errors.New("invalid syntax")
	EPlatform = errors.New("invalid platform")
)

// Parse is exported to parse a string into a PackageId struct. This function
// should only be used on user inputs. Therefore, It does NOT return an
// error but instead invokes a panic if the syntax is invalid.
func Parse(s string) (p *syntaxtypes.PackageId) {
	s = sanitize(s)
	p = &syntaxtypes.PackageId{}
	var err error
	p.Platform, p.Name, p.Version, err = parseOperatorAt(s)
	if err != nil {
		if errors.Is(err, ESyntax) {
			panic(err)
		} else {
			log.Fatal(err)
		}
	}
	return
}

// parseOperatorAt is called first since '@' operator always occur after '/' (equivalent
// to a lower priority).
func parseOperatorAt(s string) (
	pl syntaxtypes.Platform,
	n syntaxtypes.PackageName,
	v syntaxtypes.PackageVersion,
	err error,
) {
	split := strings.Split(s, "@")

	pl, n, err = parseOperatorSlash(split[0])
	if err != nil {
		return "", "", "", ESyntax
	}

	if len(split) == 1 {
		v = syntaxtypes.AllVersion
	} else if len(split) == 2 {
		v = syntaxtypes.PackageVersion(split[1])
	} else {
		return "", "", "", ESyntax
	}

	return
}

func parseOperatorSlash(s string) (
	pl syntaxtypes.Platform,
	n syntaxtypes.PackageName,
	err error,
) {
	split := strings.Split(s, "/")

	if len(split) == 1 {
		pl = syntaxtypes.AllPlatform
		n = syntaxtypes.PackageName(split[0])
		if syntaxtypes.Platform(n).Valid() {
			// Remember, all platforms are also valid packages under themselves.
			// This literal is for users to specify the platform itself. See the
			// docs for syntaxtypes.Platform for more information.
			pl = syntaxtypes.Platform(n)
			n = syntaxtypes.PackageName(pl)
		}
	} else if len(split) == 2 {
		pl = syntaxtypes.Platform(split[0])
		if !pl.Valid() {
			return "", "", EPlatform
		}
		n = syntaxtypes.PackageName(split[1])
	} else {
		return "", "", ESyntax
	}

	return
}
