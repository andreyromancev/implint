package implint

import (
	"fmt"
)

type ErrNoLevel string

func (e ErrNoLevel) Error() string {
	const t = `No level for path:
%s`
	return fmt.Sprintf(t, string(e))
}

type ErrBadImport struct {
	from, imp string
	fromLevel, impLevel int
	file string
}

func (e ErrBadImport) Error() string {
	const t = `Bad import:
level %d %s 
imports level %d %s
in file %s`
	return fmt.Sprintf(t, e.fromLevel, e.from, e.impLevel, e.imp, e.file)
}
