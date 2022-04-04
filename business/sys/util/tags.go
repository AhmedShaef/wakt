package util

// Add bulk tags to a tagset.
func Add(oldtags, tags []string) []string {
	newTags := []string{}
	for _, tag := range tags {
		newTags = append(oldtags, tag)
	}
	return newTags
}

// Remove bulk tags from a tagset.
func Remove(oldtags, tags []string) []string {
	newTags := []string{}
	for _, tag := range tags {
		for _, oldtag := range oldtags {
			if oldtag != tag {
				newTags = append(newTags, tag)
			}
		}
	}
	return newTags
}
