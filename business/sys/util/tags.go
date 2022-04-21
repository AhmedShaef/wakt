// Package util provides some common functions used in the system.
package util

// Add bulk tags to a tagSet.
func Add(oldTag, tags []string) []string {
	newTags := []string{}
	for _, tag := range tags {
		newTags = append(oldTag, tag)
	}
	return newTags
}

// Remove bulk tags from a tagSet.
func Remove(oldTag, tags []string) []string {
	newTags := []string{}
	for _, tag := range tags {
		for _, oldTag := range oldTag {
			if oldTag != tag {
				newTags = append(newTags, tag)
			}
		}
	}
	return newTags
}
