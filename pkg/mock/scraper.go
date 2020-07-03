package mock

import "github.com/oslokommune/okctl/pkg/credentials/aws/scrape"

type scraper struct {
	ScrapeFn func(string, string, string) (string, error)
}

// Scrape invokes the mocked scrape function
func (s *scraper) Scrape(username, password, mfaToken string) (string, error) {
	return s.ScrapeFn(username, password, mfaToken)
}

// NewGoodScraper returns a scraper that will succeed
func NewGoodScraper() scrape.Scraper {
	return &scraper{
		ScrapeFn: func(string, string, string) (string, error) {
			return "SAMLIsAllGood", nil
		},
	}
}
