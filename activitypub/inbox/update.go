package inbox

import (
	"context"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/activitypub/resolvers"
	log "github.com/sirupsen/logrus"
)

func handleUpdateRequest(c context.Context, activity vocab.ActivityStreamsUpdate) error {
	// We only care about update events to followers.
	if !activity.GetActivityStreamsObject().At(0).IsActivityStreamsPerson() {
		return nil
	}

	actor, err := resolvers.GetResolvedPersonFromActor(activity.GetActivityStreamsActor())
	if err != nil {
		log.Errorln(err)
		return err
	}

	iri := actor.GetJSONLDId()
	inbox := actor.GetActivityStreamsInbox().GetIRI()
	name := actor.GetActivityStreamsName().At(0).GetXMLSchemaString()
	image := actor.GetActivityStreamsIcon().At(0).GetActivityStreamsImage().GetActivityStreamsUrl().At(0).GetIRI()
	fullUsername := apmodels.GetFullUsernameFromPerson(actor)

	return persistence.UpdateFollower(iri.Get().String(), inbox.String(), name, fullUsername, image.String())
}
