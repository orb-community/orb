package types

import (
	"fmt"
	"strconv"
	"strings"
)

// todo change to the dependency semantic version
type SemVer struct {
	Major int
	Minor int
	Patch int
}

func NewSemVerFromString(version string) SemVer {
	vPieces := strings.Split(version, ".")
	if len(vPieces) != 3 {
		return SemVer{}
	}
	major, _ := strconv.Atoi(vPieces[0])
	minor, _ := strconv.Atoi(vPieces[1])
	patch, _ := strconv.Atoi(vPieces[2])
	return SemVer{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}

func (v1 SemVer) ToString() string {
	return fmt.Sprintf("%d.%d.%d", v1.Major, v1.Minor, v1.Patch)
}

func (v1 SemVer) IsEqualTo(v2 SemVer) bool {
	if v1.Major == v2.Major && v1.Minor == v2.Minor && v1.Patch == v2.Patch {
		return true
	}
	return false
}

func (v1 SemVer) IsNewerThan(v2 SemVer) bool {
	if v1.Major > v2.Major {
		return true
	} else if v1.Major == v2.Major {
		if v1.Minor > v2.Minor {
			return true
		} else if v1.Minor == v2.Minor {
			if v1.Patch > v2.Patch {
				return true
			}
		}
	}
	return false
}
