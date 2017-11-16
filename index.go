package main

import (
	"os"

	"github.com/blevesearch/bleve"

	log "github.com/sirupsen/logrus"
)

func unIndex(uid string) error {
	if err := idx.Delete(uid); err != nil {
		log.Errorf("UnIndexed uid=%s success=false", uid)
		return err
	}
	log.Infof("UnIndexed uid=%s success=true", uid)
	return nil
}

func getIndex(path string) (bleve.Index, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		mapping := bleve.NewIndexMapping()
		idx, err := bleve.New(path, mapping)
		if err != nil {
			return idx, err
		}
		idx.Close()
	}
	return bleve.OpenUsing(path, map[string]interface{}{
		"read_only": *secondary,
	})
}

func indexNote(note Note) error {
	err := idx.Index(note.UID, note)
	if err != nil {
		log.Errorf("Indexed uid=%s success=false", note.UID)
		return err
	}
	log.Infof("Indexed uid=%s success=true", note.UID)
	return err
}

func searchIndex(q string) ([]string, error) {
	query := bleve.NewQueryStringQuery(q)
	in := bleve.NewSearchRequest(query)
	out, err := idx.Search(in)
	if err != nil {
		return nil, err
	}
	uids := []string{}
	for _, hit := range out.Hits {
		uids = append(uids, string(hit.IndexInternalID))
	}
	return uids, nil
}

func reIndex(b Backend) error {
	for _, note := range b.list() {
		content, _ := getContentByUID(b, note.UID, "")
		note.Content = content
		err := indexNote(note)
		if err != nil {
			return err
		}
	}
	return nil
}
