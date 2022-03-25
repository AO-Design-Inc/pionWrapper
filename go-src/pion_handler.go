// a handler for pion
package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/frame"
	"github.com/pion/mediadevices/pkg/prop"
	"github.com/pion/webrtc/v3"

	"github.com/pion/mediadevices/pkg/codec/x264"

	_ "github.com/pion/mediadevices/pkg/driver/screen"

	"encoding/json"
	"fmt"
)

type JSONString []byte

var peerConnection *webrtc.PeerConnection

//export SpawnConnection
func SpawnConnection(iceValues JSONString) *C.char {
	var iceServers []webrtc.ICEServer
	if err := json.Unmarshal(iceValues, &iceServers); err != nil {
		panic(err)
	}

	config := webrtc.Configuration{
		ICEServers: iceServers,
	}

	x264Params, err := x264.NewParams()
	if err != nil {
		panic(err)
	}

	codecSelector := mediadevices.NewCodecSelector(
		mediadevices.WithVideoEncoders(&x264Params),
	)

	mediaEngine := webrtc.MediaEngine{}
	codecSelector.Populate(&mediaEngine)
	api := webrtc.NewAPI(webrtc.WithMediaEngine(&mediaEngine))
	peerConnection, err = api.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if cErr := peerConnection.Close(); cErr != nil {
			panic(cErr)
		}
	}()

	stream, err := mediadevices.GetDisplayMedia(mediadevices.MediaStreamConstraints{
		Video: func(constraint *mediadevices.MediaTrackConstraints) {
			constraint.FrameFormat = prop.FrameFormat(frame.FormatI420)
			constraint.FrameRate = prop.Float(60)
		},
		Codec: codecSelector,
	})
	if err != nil {
		panic(err)
	}

	for _, track := range stream.GetTracks() {
		track.OnEnded(func(err error) {
			fmt.Printf("Track (ID: %s) ended with error: %v\n",
				track.ID(), err)
		})

		_, err = peerConnection.AddTransceiverFromTrack(track,
			webrtc.RtpTransceiverInit{
				Direction: webrtc.RTPTransceiverDirectionSendonly,
			},
		)
		if err != nil {
			panic(err)
		}
	}

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState)
	})

	track := stream.GetVideoTracks()[0].(*mediadevices.VideoTrack)
	defer track.Close()

	const mtu = 1000
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	if err = peerConnection.SetLocalDescription(offer); err != nil {
		panic(err)
	}

	offerString, err := json.Marshal(offer)
	cOfferString := C.CString(string(offerString))
	//defer C.free(cOfferString)
	return cOfferString
}

//export SetRemoteDescription
func SetRemoteDescription(remoteDescString JSONString) bool {
	var desc webrtc.SessionDescription
	err := json.Unmarshal(remoteDescString, &desc)
	if err != nil {
		return false
	}
	peerConnection.SetRemoteDescription(desc)
	return true
}


func main(){}
