package models

type GoogleBookmarks struct {
	Checksum string              `json:"checksum,omitempty"`
	Roots    GoogleBookmarksRoot `json:"roots,omitempty"`
	Version  int                 `json:"version,omitempty"`
}

type GoogleBookmarksRoot struct {
	BookmarkBar GoogleBookmarkItem `json:"bookmark_bar,omitempty"`
	Other       GoogleBookmarkItem `json:"other,omitempty"`
	Synced      GoogleBookmarkItem `json:"synced,omitempty"`
}

type GoogleBookmarkItem struct {
	Children     []GoogleBookmarkItem `json:"children,omitempty"` //only for folder
	DateAdded    string               `json:"date_added,omitempty"`
	DateModified string               `json:"date_modified,omitempty"` //only for folder
	Guid         string               `json:"guid,omitempty"`
	Id           string               `json:"id,omitempty"`
	Name         string               `json:"name,omitempty"`
	Type         string               `json:"type,omitempty"` // folder or url
	Url          string               `json:"url,omitempty"`  //only for url
}

type GoogleBookmarkFlat struct {
	DateAdded string `json:"date_added,omitempty"`
	Folder    string `json:"folder,omitempty"`
	Guid      string `json:"guid,omitempty"`
	Id        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Type      string `json:"type,omitempty"`
	Url       string `json:"url,omitempty"`
}

func (gb *GoogleBookmarks) GetFlatList() []GoogleBookmarkFlat {
	bookmarks := make([]GoogleBookmarkFlat, 0)
	for _, value := range gb.Roots.BookmarkBar.Children {
		if value.Type == "folder" {
			bookmarks = append(bookmarks, value.getChildren(value.Name)...)
		} else {
			bookmarks = append(bookmarks, value.toFlat("main"))
		}
	}
	for _, value := range gb.Roots.Other.Children {
		if value.Type == "folder" {
			bookmarks = append(bookmarks, value.getChildren(value.Name)...)
		} else {
			bookmarks = append(bookmarks, value.toFlat("other"))
		}
	}
	return bookmarks
}

func (gbi *GoogleBookmarkItem) getChildren(folder string) []GoogleBookmarkFlat {
	bookmarks := make([]GoogleBookmarkFlat, 0)
	for _, value := range gbi.Children {
		bookmarks = append(bookmarks, value.toFlat(gbi.Name))
	}
	return bookmarks
}

func (gbi *GoogleBookmarkItem) toFlat(folder string) GoogleBookmarkFlat {
	return GoogleBookmarkFlat{
		DateAdded: gbi.DateAdded,
		Folder:    folder,
		Guid:      gbi.Guid,
		Id:        gbi.Id,
		Name:      gbi.Name,
		Type:      gbi.Type,
		Url:       gbi.Url,
	}
}
