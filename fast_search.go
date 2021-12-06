package goinsta

func (search *Search) FastSearchUser(query string) (*SearchResult, error) {
	return search.fastSearch(query, search.user)
}

func (search *Search) FastSearchHashtag(query string) (*SearchResult, error) {
	return search.fastSearch(query, search.tags)
}

func (search *Search) fastSearch(query string, fn func(string) (*SearchResult, error)) (*SearchResult, error) {
	result, err := fn(query)
	if err != nil {
		return nil, err
	}
	// If the query is a username, and in the top 10, return
	if len(result.Results) >= 10 {
		for _, r := range result.Results[:10] {
			if r.User != nil && r.User.Username == query {
				return result, nil
			}
		}
	}

	return result, nil
}
