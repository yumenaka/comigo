package model

import "time"

type BookMarks []BookMark

func (s *BookMarks) GetLastReadPage() int {
	for _, mark := range *s {
		if mark.Type == AutoMark {
			return mark.PageIndex
		}
	}
	return 0
}

func (s *BookMarks) GetLastReadTime() time.Time {
	for _, mark := range *s {
		if mark.Type == AutoMark {
			return mark.UpdatedAt
		}
	}
	return time.Time{}
}
