package flickr

func (c *PhotosClient) Favs(userId string) (*PhotoList, error) {
	response, err := c.Request("favorites.getPublicList", Params{
		"user_id": userId,
		// "extras":  "license",
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}
