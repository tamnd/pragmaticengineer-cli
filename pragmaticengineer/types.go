package pragmaticengineer

type Post struct {
	Rank     int    `json:"rank"     csv:"rank"     tsv:"rank"`
	Date     string `json:"date"     csv:"date"     tsv:"date"`
	Audience string `json:"audience" csv:"audience" tsv:"audience"`
	Title    string `json:"title"    csv:"title"    tsv:"title"`
	Subtitle string `json:"subtitle" csv:"subtitle" tsv:"subtitle"`
	URL      string `json:"url"      csv:"url"      tsv:"url"`
}
