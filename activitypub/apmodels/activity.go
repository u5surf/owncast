package apmodels

import (
	"net/url"
	"time"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
)

// PrivacyAudience represents the audience for an activity.
type PrivacyAudience = string

const (
	// PUBLIC is an audience meaning anybody can view the item.
	PUBLIC PrivacyAudience = "https://www.w3.org/ns/activitystreams#Public"
)

// MakeCreateActivity will return a new Create activity with the provided ID.
func MakeCreateActivity(activityID *url.URL) vocab.ActivityStreamsCreate {
	activity := streams.NewActivityStreamsCreate()
	id := streams.NewJSONLDIdProperty()
	id.Set(activityID)
	activity.SetJSONLDId(id)

	public, _ := url.Parse(PUBLIC)
	to := streams.NewActivityStreamsToProperty()
	to.AppendIRI(public)
	activity.SetActivityStreamsTo(to)

	return activity
}

// MakeUpdateActivity will return a new Update activity with the provided aID.
func MakeUpdateActivity(activityID *url.URL) vocab.ActivityStreamsUpdate {
	activity := streams.NewActivityStreamsUpdate()
	id := streams.NewJSONLDIdProperty()
	id.Set(activityID)
	activity.SetJSONLDId(id)

	public, _ := url.Parse(PUBLIC)
	to := streams.NewActivityStreamsToProperty()
	to.AppendIRI(public)
	activity.SetActivityStreamsTo(to)

	return activity
}

// MakeNote will return a new Note object.
func MakeNote(text string, noteIRI *url.URL, attributedToIRI *url.URL) vocab.ActivityStreamsNote {
	note := streams.NewActivityStreamsNote()
	content := streams.NewActivityStreamsContentProperty()
	content.AppendXMLSchemaString(text)
	note.SetActivityStreamsContent(content)
	id := streams.NewJSONLDIdProperty()
	id.Set(noteIRI)
	note.SetJSONLDId(id)

	published := streams.NewActivityStreamsPublishedProperty()
	published.Set(time.Now())
	note.SetActivityStreamsPublished(published)

	attributedTo := attributedToIRI
	attr := streams.NewActivityStreamsAttributedToProperty()
	attr.AppendIRI(attributedTo)
	note.SetActivityStreamsAttributedTo(attr)

	public, _ := url.Parse(PUBLIC)
	to := streams.NewActivityStreamsToProperty()
	to.AppendIRI(public)
	note.SetActivityStreamsTo(to)

	return note
}
