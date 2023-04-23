package database

type Content struct {
	Text string `json:"text"`

	// Image image.Image
}

func (c *Content) IsEmpty() bool {
	return c.Text == ""
}
