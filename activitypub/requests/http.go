package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/crypto"
)

func WriteStreamResponse(item vocab.Type, w http.ResponseWriter, publicKey apmodels.PublicKey) error {
	var jsonmap map[string]interface{}
	jsonmap, _ = streams.Serialize(item)
	b, err := json.Marshal(jsonmap)
	if err != nil {
		return err
	}

	return WriteResponse(b, w, publicKey)
}

func WritePayloadResponse(payload interface{}, w http.ResponseWriter, publicKey apmodels.PublicKey) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return WriteResponse(b, w, publicKey)
}

func WriteResponse(payload []byte, w http.ResponseWriter, publicKey apmodels.PublicKey) error {
	w.Header().Set("Content-Type", "application/json")

	if err := crypto.SignResponse(w, payload, publicKey); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return err
	}

	fmt.Println(string(payload))
	if _, err := w.Write(payload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	fmt.Println(string(payload))

	return nil
}

func PostSignedRequest(payload []byte, url *url.URL, fromActorIRI *url.URL) ([]byte, error) {
	fmt.Println("Sending", string(payload), "to", url)

	req, _ := http.NewRequest("POST", url.String(), bytes.NewBuffer(payload))
	if err := crypto.SignRequest(req, payload, fromActorIRI); err != nil {
		fmt.Println("error signing request:", err)
		return nil, err
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, _ := ioutil.ReadAll(response.Body)

	fmt.Println("Response: ", response.StatusCode, string(body))
	return body, nil
}
