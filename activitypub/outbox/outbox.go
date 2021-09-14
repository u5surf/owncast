package outbox

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/activitypub/requests"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
)

// SendLive will send all followers the message saying you started a live stream.
func SendLive() error {
	textContent := data.GetFederationGoLiveMessage()
	textContent = utils.RenderSimpleMarkdown(textContent)

	localActor := apmodels.MakeLocalIRIForAccount(data.GetDefaultFederationUsername())
	noteID := shortid.MustGenerate()
	noteIRI := apmodels.MakeLocalIRIForResource(noteID)
	id := shortid.MustGenerate()
	activity := apmodels.CreateCreateActivity(id, localActor)
	object := streams.NewActivityStreamsObjectProperty()
	activity.SetActivityStreamsObject(object)

	tagStrings := []string{}
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")

	tagProp := streams.NewActivityStreamsTagProperty()
	for _, tagString := range data.GetServerMetadataTags() {
		tagWithoutSpecialCharacters := reg.ReplaceAllString(tagString, "")
		hashtag := apmodels.MakeHashtag(tagWithoutSpecialCharacters)
		tagProp.AppendTootHashtag(hashtag)

		// TODO: Do we want to display tags or just assign them?
		tagString := fmt.Sprintf("<a class=\"hashtag\" href=\"https://directory.owncast.online/tags/%s\">#%s</a>", tagWithoutSpecialCharacters, tagWithoutSpecialCharacters)
		tagStrings = append(tagStrings, tagString)
	}

	// Manually add Owncast hashtag if it doesn't already exist so it shows up
	// in Owncast search results.
	// We can remove this down the road, but it'll be nice for now.
	if _, exists := utils.FindInSlice(tagStrings, "owncast"); !exists {
		hashtag := apmodels.MakeHashtag("owncast")
		tagProp.AppendTootHashtag(hashtag)
	}

	tagsString := strings.Join(tagStrings, " ")

	var streamTitle string
	if title := data.GetStreamTitle(); title != "" {
		streamTitle = fmt.Sprintf("<p>%s</p>", title)
	}
	textContent = fmt.Sprintf("<p>%s</p><a href=\"%s\">%s</a>%s<p>%s</p>", textContent, data.GetServerURL(), data.GetServerURL(), streamTitle, tagsString)

	log.Println(textContent)

	note := apmodels.MakeNote(textContent, noteIRI, localActor)
	note.SetActivityStreamsTag(tagProp)
	object.AppendActivityStreamsNote(note)

	// Attach an image along with the Federated message.
	previewURL, err := url.Parse(data.GetServerURL())
	if err == nil {
		var imageToAttach string
		previewGif := filepath.Join(config.WebRoot, "preview.gif")
		thumbnailJpg := filepath.Join(config.WebRoot, "thumbnail.jpg")

		if utils.DoesFileExists(previewGif) {
			imageToAttach = "preview.gif"
		} else if utils.DoesFileExists(thumbnailJpg) {
			imageToAttach = "thumbnail.jpg"
		}
		if imageToAttach != "" {
			previewURL.Path = imageToAttach
			apmodels.AddImageAttachmentToNote(note, previewURL.String())
		}
	}

	b, err := apmodels.Serialize(activity)
	if err != nil {
		log.Errorln("unable to serialize go live message activity", err)
		return errors.New("unable to serialize go live message activity " + err.Error())
	}

	if err := SendToFollowers(b); err != nil {
		return err
	}

	if err := Add(activity, id); err != nil {
		return err
	}
	if err := Add(note, noteID); err != nil {
		return err
	}

	return nil
}

// SendPublicMessage will send a public message to all followers.
func SendPublicMessage(textContent string) error {
	localActor := apmodels.MakeLocalIRIForAccount(data.GetDefaultFederationUsername())
	id := shortid.MustGenerate()
	activity := apmodels.CreateMessageActivity(id, textContent, localActor)
	message := activity.GetActivityStreamsObject().Begin().GetActivityStreamsNote()

	b, err := apmodels.Serialize(activity)
	if err != nil {
		return errors.New("unable to serialize the send public message activity")
	}
	if err := SendToFollowers(b); err != nil {
		return err
	}

	return Add(message, message.GetJSONLDId().Get().String())
}

// SendToFollowers will send an arbitrary payload to all follower inboxes.
func SendToFollowers(payload []byte) error {
	localActor := apmodels.MakeLocalIRIForAccount(data.GetDefaultFederationUsername())

	followers, err := persistence.GetFederationFollowers()
	if err != nil {
		log.Errorln("unable to fetch followers to send to", err)
		return errors.New("unable to fetch followers to send payload to")
	}

	for _, follower := range followers {
		inbox, _ := url.Parse(follower.Inbox)
		if _, err := requests.PostSignedRequest(payload, inbox, localActor); err != nil {
			log.Errorln("unable to send to follower inbox", follower.Inbox, err)
			return errors.New("unable to send to follower inbox: " + follower.Inbox)
		}
	}
	return nil
}

// UpdateFollowersWithAccountUpdates will send an update to all followers alerting of a profile update.
func UpdateFollowersWithAccountUpdates() error {
	log.Println("Updating followers with new actor details")

	id := shortid.MustGenerate()
	objectID := apmodels.MakeLocalIRIForResource(id)
	activity := apmodels.MakeUpdateActivity(objectID)

	actor := streams.NewActivityStreamsPerson()
	actorID := apmodels.MakeLocalIRIForAccount(data.GetDefaultFederationUsername())
	actorIDProperty := streams.NewJSONLDIdProperty()
	actorIDProperty.Set(actorID)
	actor.SetJSONLDId(actorIDProperty)

	actorProperty := streams.NewActivityStreamsActorProperty()
	actorProperty.AppendActivityStreamsPerson(actor)
	activity.SetActivityStreamsActor(actorProperty)

	obj := streams.NewActivityStreamsObjectProperty()
	obj.AppendIRI(actorID)
	activity.SetActivityStreamsObject(obj)

	b, err := apmodels.Serialize(activity)
	if err != nil {
		log.Errorln("unable to serialize send update actor activity", err)
		return errors.New("unable to serialize send update actor activity")
	}
	return SendToFollowers(b)
}

// Add will save an ActivityPub object to the datastore.
func Add(item vocab.Type, id string) error {
	iri := "/" + item.GetJSONLDId().GetIRI().Path
	typeString := item.GetTypeName()

	if iri == "" {
		log.Errorln("Unable to get iri from item")
		return errors.New("Unable to get iri from item " + id)
	}

	b, err := apmodels.Serialize(item)
	if err != nil {
		log.Errorln("unable to serialize model when saving to outbox", err)
		return err
	}

	return persistence.AddToOutbox(id, iri, b, typeString)
}

// Get will return the outbox.
func Get() vocab.ActivityStreamsOrderedCollectionPage {
	orderedCollection, _ := persistence.GetOutbox()
	return orderedCollection
}
