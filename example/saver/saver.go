package saver

import (
	// This is ok.
	_ "github.com/andreyromancev/implint/example/db"
	// This should fail. Saver should not have access to this API.
	_ "github.com/andreyromancev/implint/example/fetcher"
	// This is ok. Nolint skips the error.
	_ "github.com/andreyromancev/implint/example" // nolint
)
