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


func peerConnector(config *webrtc.Configuration, recvSdp chan *C.char) {
	x264Params, err := x264.NewParams()
	if err != nil {
		panic(err)
	}
	x264Params.Preset = x264.PresetMedium
	x264Params.BitRate = 1_000_000

	codecSelector := mediadevices.NewCodecSelector(
		mediadevices.WithVideoEncoders(&x264Params),
	)

	mediaEngine := webrtc.MediaEngine{}
	codecSelector.Populate(&mediaEngine)
	api := webrtc.NewAPI(webrtc.WithMediaEngine(&mediaEngine))
	peerConnection, err = api.NewPeerConnection(*config)
	if err != nil {
		panic(err)
	}

	stream, err := mediadevices.GetDisplayMedia(mediadevices.MediaStreamConstraints{
		Video: func(constraint *mediadevices.MediaTrackConstraints) {
      constraint.Width = prop.Int(640)
      constraint.Height = prop.Int(480)
			constraint.FrameFormat = prop.FrameFormat(frame.FormatI420)
			constraint.FrameRate = prop.Float(60)
		},
		Codec: codecSelector,
	})

	for _, track := range stream.GetTracks() {
    fmt.Printf("%v\n", track)
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
	if err != nil {
		panic(err)
	}

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState)
	})
	const mtu = 1000
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)
	if err = peerConnection.SetLocalDescription(offer); err != nil {
		panic(err)
	}

	<-gatherComplete
	offerString, err := json.Marshal(*peerConnection.LocalDescription())
	cOfferString := C.CString(string(offerString))
	recvSdp <- cOfferString
  select {}
}

//export SpawnConnection
func SpawnConnection(iceValues JSONString) *C.char {
	sdpRecv := make(chan *C.char, 1)
	var iceServers []webrtc.ICEServer
	if err := json.Unmarshal(iceValues, &iceServers); err != nil {
		panic(err)
	}

	config := webrtc.Configuration{
		ICEServers: iceServers,
	}

	go peerConnector(&config, sdpRecv)

	
	return(<-sdpRecv)
}

//export SetRemoteDescription
func SetRemoteDescription(remoteDescString JSONString) bool {
	var desc webrtc.SessionDescription
	if err := json.Unmarshal(remoteDescString, &desc); err != nil {
		return false
	}
  go remoteSetter(&desc)
	return true
}

func remoteSetter(desc *webrtc.SessionDescription) {
	if err := peerConnection.SetRemoteDescription(*desc); err != nil {
		panic(err)
	}

  select {}
}


func main() {}
