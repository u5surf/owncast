package persistence

import (
	"fmt"
	"net/url"

	"github.com/owncast/owncast/activitypub/models"
	"github.com/owncast/owncast/core/data"
	log "github.com/sirupsen/logrus"
)

var _datastore *data.Datastore

func Setup(datastore *data.Datastore) {
	_datastore = datastore
	createFederationFollowersTable()
}

func AddFollow(follow models.ActivityPubActor) error {
	fmt.Println("Saving", follow.ActorIri, "as a follower.")
	return createFollow(follow.ActorIri, follow.Inbox)
}

func RemoveFollow(unfollow models.ActivityPubActor) error {
	fmt.Println("Removing", unfollow.ActorIri, "as a follower.")
	return removeFollow(unfollow.ActorIri)
}

func GetFederationFollowers() ([]models.ActivityPubActor, error) {
	followers := make([]models.ActivityPubActor, 0)

	var query = "SELECT iri, inbox FROM ap_followers"

	rows, err := _datastore.DB.Query(query)
	if err != nil {
		return followers, err
	}
	defer rows.Close()

	for rows.Next() {
		var iriString string
		var inboxString string

		if err := rows.Scan(&iriString, &inboxString); err != nil {
			log.Error("There is a problem reading the database.", err)
			return followers, err
		}

		iri, _ := url.Parse(iriString)
		inbox, _ := url.Parse(inboxString)

		singleFollower := models.ActivityPubActor{
			ActorIri: iri,
			Inbox:    inbox,
		}

		followers = append(followers, singleFollower)
	}

	if err := rows.Err(); err != nil {
		return followers, err
	}

	return followers, nil
}

func createFederationFollowersTable() {
	log.Traceln("Creating federation followers table...")

	createTableSQL := `CREATE TABLE IF NOT EXISTS ap_followers (
		"iri" TEXT NOT NULL,
		"inbox" TEXT NOT NULL,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (iri));
		CREATE INDEX iri ON ap_followers (iri);
	);`

	stmt, err := _datastore.DB.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		log.Warnln(err)
	}
}

func createFollow(actor *url.URL, inbox *url.URL) error {
	_datastore.DbLock.Lock()
	defer _datastore.DbLock.Unlock()

	tx, err := _datastore.DB.Begin()
	if err != nil {
		log.Debugln(err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.Prepare("INSERT INTO ap_followers(iri, inbox) values(?, ?)")

	if err != nil {
		log.Debugln(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(actor.String(), inbox.String())
	if err != nil {
		log.Errorln("error creating new federation follow", err)
	}

	return tx.Commit()
}

func removeFollow(actor *url.URL) error {
	_datastore.DbLock.Lock()
	defer _datastore.DbLock.Unlock()

	tx, err := _datastore.DB.Begin()
	if err != nil {
		log.Debugln(err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.Prepare("DELETE FROM ap_followers WHERE iri IS ?")

	if err != nil {
		log.Debugln(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(actor.String())
	if err != nil {
		log.Errorln("error removing federation follow", err)
	}

	return tx.Commit()
}
