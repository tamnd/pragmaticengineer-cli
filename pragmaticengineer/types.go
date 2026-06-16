package pragmaticengineer

type Post struct {
	Rank     int    `json:"rank"     csv:"rank"     tsv:"rank"`
	Date     string `json:"date"     csv:"date"     tsv:"date"`
	Audience string `json:"audience" csv:"audience" tsv:"audience"`
	Title    string `json:"title"    csv:"title"    tsv:"title"`
	Subtitle string `json:"subtitle" csv:"subtitle" tsv:"subtitle"`
	URL      string `json:"url"      csv:"url"      tsv:"url"`
}

// Info holds aggregate statistics about the newsletter.
type Info struct {
	TotalPosts  int    `json:"total_posts"`
	FreePosts   int    `json:"free_posts"`
	PaidPosts   int    `json:"paid_posts"`
	OldestPost  string `json:"oldest_post"`
	LatestPost  string `json:"latest_post"`
	NewsletterURL string `json:"newsletter_url"`
}
