package models

import (
	"time"
)


type FileInfo struct {
	Name 		 string	`json:"name"`
	OriginalName string `json:"original_name, omitempty"`
	Dir 		 string `json:"dir"`
	Path		 string `json:"path"`
	Extension 	 string `json:"extension"`
	Size     	 int64  `json:"size"`
	ModTime  	 time.Time `json:"mod_time"`
	IsDir        bool   `json:"is_dir"`
	MimeType     string `json:"mime"`
}


type ListingInfo struct {
	Files []FileInfo   `json:"files"`
}



