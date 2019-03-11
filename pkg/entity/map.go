package entity

// Map contains a geolocation's information.
type Map struct {
	ID                string `json:"id"`
	Desctription      string `json:"description"`
	MatchedSubstrings []struct {
		Length int `json:"length"`
		Offset int `json:"secondary_text"`
	}
	PlaceID              string `json:"place_id"`
	Reference            string `json:"reference"`
	StructuredFormatting struct {
		MainText                  string `json:"main_text"`
		MainTextMatchedSubstrings []struct {
			Length int `json:"length"`
			Offset int `json:"secondary_text"`
		} `json:"main_text_matched_substrings"`
		SecondaryText string `json:"secondary_text"`
	} `json:"structured_formatting"`
	Terms []struct {
		Offset int    `json:"offset"`
		Value  string `json:"value"`
	} `json:"terms"`
	Types []string `json:"types"`
}

// Validate validates that the map's required fields are filled out correctly.
func (t *Map) Validate() error {
	return nil
}
