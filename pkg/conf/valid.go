package conf

import "errors"

func (c *Conf) Validate() error {
	pCount := len(c.Paths.Movies)
	pCount += len(c.Paths.Music)
	pCount += len(c.Paths.Series)
	if pCount == 0 {
		return errors.New("No search paths configured")
	}

	if (c.Backend.ArangoDB == nil) && (c.Backend.OpenSearch == nil) {
		return errors.New("No backend configured")
	}

	if c.Backend.ArangoDB != nil {
		if c.Backend.ArangoDB.Auth.User == "" {
			return errors.New("ArangoDB configuration missing 'auth.user'")
		}
		if c.Backend.ArangoDB.Auth.Pass == "" {
			return errors.New("ArangoDB configuration missing 'auth.pass'")
		}
	}

	if c.Backend.OpenSearch != nil {
		if c.Backend.OpenSearch.Auth.User == "" {
			return errors.New("OpenSearch configuration missing 'auth.user'")
		}
		if c.Backend.OpenSearch.Auth.Pass == "" {
			return errors.New("OpenSearch configuration missing 'auth.pass'")
		}
	}

	return nil
}
