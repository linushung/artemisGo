package postgres

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func (rdb *RDB) SelectArticleById(id uuid.UUID) (Article, error) {
	a := Article{}
	statement := `SELECT * FROM article WHERE id = ?;`
	statement = rdb.Poolx.Rebind(statement)

	if err := rdb.Poolx.Get(&a, statement, id); err != nil {
		log.Errorf("***** [POSTGRES:%s][FAIL] ***** Cannot execute SELECT operation:: %v", "SelectArticleById", err)
		return a, err
	}

	return a, nil
}

func (rdb *RDB) CreateArticle(id uuid.UUID, article Article) (Article, error) {
	articleStmt := `INSERT INTO article (id, slug, title, description, body) VALUES (?,?,?,?,?);`
	articleStmt = rdb.Poolx.Rebind(articleStmt)

	_, err := rdb.Poolx.Exec(articleStmt,
		id,
		article.Slug,
		article.Title,
		article.Description,
		article.Body,
	)
	if err != nil {
		log.Errorf("***** [POSTGRES:%s][FAIL] ***** Cannot execute INSERT operation:: %v", "CreateArticle", err)
		return Article{}, err
	}

	a, err := rdb.SelectArticleById(id)
	if err != nil {
		log.Errorf("***** [POSTGRES:%s][FAIL] ***** Cannot execute SELECT operation:: %v", "CreateArticle", err)
		return Article{}, err
	}

	return a, nil
}

func (rdb *RDB) TagArticle(id int64, tags []string) error {
	for _, t := range tags {
		tagStmt := `INSERT INTO tag (id, tag) VALUES (?,?);`
		tagStmt = rdb.Poolx.Rebind(tagStmt)

		if _, err := rdb.Poolx.Exec(tagStmt, id, t); err != nil {
			log.Errorf("***** [POSTGRES:%s][FAIL] ***** Cannot execute INSERT operation:: %v", "TagArticle", err)
			return err
		}
	}

	return nil
}
