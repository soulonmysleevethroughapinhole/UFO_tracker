package web

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v2"
	"github.com/pion/webrtc/v2/pkg/media"
	"github.com/soulonmysleevethroughapinhole/UFO_tracker/agent"
)

//
// Testing
//

const (
	rtcpPLIInterval = time.Second * 3
	compress        = false
)

var peerConnectionConfig = webrtc.Configuration{
	ICEServers: []webrtc.ICEServer{
		{
			URLs: []string{"stun:stun.l.google.com:19302"},
		},
	},
}

func saveToDisk(i media.Writer, track *webrtc.Track) {
	defer func() {
		if err := i.Close(); err != nil {
			panic(err)
		}
	}()

	for {
		rtpPacket, err := track.ReadRTP()
		if err != nil {
			panic(err)
		}
		if err := i.WriteRTP(rtpPacket); err != nil {
			panic(err)
		}
	}
}

//not quite sure how this works, maybe change and see what happens
var (
	incomingClients = make(chan *agent.Client, 10)
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	EnableCompression: true,
}

type WebInterface struct {
	Clients     map[int]*agent.Client
	ClientsLock *sync.Mutex

	Channels     map[int]*agent.Channel
	ChannelsLock *sync.Mutex

	Description  string
	Image        string
	Artist       string
	NumOfPeers   int
	SetToDelete  bool
	MediaWriting bool
}

func NewWebInterface(router *chi.Mux, path string) *WebInterface {
	w := WebInterface{Clients: make(map[int]*agent.Client),
		ClientsLock:  new(sync.Mutex),
		Channels:     make(map[int]*agent.Channel),
		ChannelsLock: new(sync.Mutex),
		Description:  "Description not set",
		Image:        "default.png",
		Artist:       "username",
		NumOfPeers:   0,
		SetToDelete:  false,
		MediaWriting: true,
	}

	w.ChannelsLock.Lock()
	defer w.ChannelsLock.Unlock()

	var (
		sdpChan    = make(chan string)
		answerChan = make(chan []byte)
	)

	regPath := "/api/sdp/" + path

	log.Println(regPath)

	//getting the sdp is stateless, people can keep the connection even if they are not in the chit chat room
	router.Post(regPath, func(wr http.ResponseWriter, r *http.Request) {

		//log.Println(regPath)

		body, _ := ioutil.ReadAll(r.Body)
		sdpChan <- string(body)

		//log.Println(w.NumOfPeers)

		answer := <-answerChan
		fmt.Fprint(wr, string(answer))
	})
	//this is for the chat only, and managing connections
	//router.Get(wsPath, w.webSocketHandler)
	/* This won't be needed yet */

	go func() {
		/* cracks knuckles */
		/* nothin personnel kid */

		offer := webrtc.SessionDescription{}

		DecodeBase64(<-sdpChan, &offer)

		mediaEngine := webrtc.MediaEngine{}
		err := mediaEngine.PopulateFromSDP(offer)
		if err != nil {
			panic(err)
		}
		api := webrtc.NewAPI(webrtc.WithMediaEngine(mediaEngine))

		peerConnection, err := api.NewPeerConnection(peerConnectionConfig)
		if err != nil {
			panic(err)
		}

		if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
			panic(err)
		}

		localTrackChan := make(chan *webrtc.Track)
		peerConnection.OnTrack(func(remoteTrack *webrtc.Track, receiver *webrtc.RTPReceiver) {
			go func() {
				ticker := time.NewTicker(rtcpPLIInterval)
				for range ticker.C {
					if rtcpSendErr := peerConnection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: remoteTrack.SSRC()}}); rtcpSendErr != nil {
						fmt.Println(rtcpSendErr)
					}
				}
			}()

			localTrack, newTrackErr := peerConnection.NewTrack(remoteTrack.PayloadType(), remoteTrack.SSRC(), "video", "pion")
			if newTrackErr != nil {
				panic(newTrackErr)
			}
			localTrackChan <- localTrack

			rtpBuf := make([]byte, 1400)
			for {
				i, readErr := remoteTrack.Read(rtpBuf)
				if readErr != nil {
					panic(readErr)
				}

				// ErrClosedPipe means we don't have any subscribers, this is ok if no peers have connected yet
				if _, err = localTrack.Write(rtpBuf[:i]); err != nil && err != io.ErrClosedPipe {
					panic(err)
				}
			}
		})

		err = peerConnection.SetRemoteDescription(offer)
		if err != nil {
			panic(err)
		}

		answer, err := peerConnection.CreateAnswer(nil)
		if err != nil {
			panic(err)
		}

		err = peerConnection.SetLocalDescription(answer)
		if err != nil {
			panic(err)
		}

		log.Println("answer sent to first peer")
		answerChan <- []byte(EncodeBase64(answer))

		localTrack := <-localTrackChan
		for {
			w.NumOfPeers++ //I THINK THIS SHOULD WORK ?!?!?!? CALLED EVERY TIME MAYBE?!?!?!

			recvOnlyOffer := webrtc.SessionDescription{}
			DecodeBase64(<-sdpChan, &recvOnlyOffer)

			peerConnection, err := api.NewPeerConnection(peerConnectionConfig)
			if err != nil {
				panic(err)
			}

			_, err = peerConnection.AddTrack(localTrack)
			if err != nil {
				panic(err)
			}

			err = peerConnection.SetRemoteDescription(recvOnlyOffer)
			if err != nil {
				panic(err)
			}

			answer, err := peerConnection.CreateAnswer(nil)
			if err != nil {
				panic(err)
			}

			err = peerConnection.SetLocalDescription(answer)
			if err != nil {
				panic(err)
			}

			//man I totally forgot how this works
			log.Printf("\nanswer sent to %d peer", w.NumOfPeers)
			answerChan <- []byte(EncodeBase64(answer))
		}
	}()

	return &w
}

// func (w *WebInterface) createChannels() {
// 	w.ChannelsLock.Lock()
// 	defer w.ChannelsLock.Unlock()

// 	// TODO Load channels from database
// }

// func (w *WebInterface) AddChannel(name string, topic string) {
// 	w.ChannelsLock.Lock()
// 	defer w.ChannelsLock.Unlock()

// 	//ch := agent.NewChannel(w.nextChannelID(), t)
// 	ch := agent.NewChannel(w.nextChannelID())
// 	ch.Name = name
// 	ch.Topic = topic

// 	w.Channels[ch.ID] = ch
// }

// func (w *WebInterface) nextChannelID() int {
// 	id := 0
// 	for cid := range w.Channels {
// 		if cid > id {
// 			id = cid
// 		}
// 	}

// 	return id + 1
// }

// /*yeah no idea anymore on what is happening*/
// //func (w *WebInterface) MessageRequestSDP(offerSDP []byte) ([]byte, error) {

// func (w *WebInterface) nextClientID() int {
// 	id := 1
// 	for {
// 		if _, ok := w.Clients[id]; !ok {
// 			break
// 		}

// 		id++
// 	}
// 	return id
// }

// func (w *WebInterface) sendChannelList(c *agent.Client) {
// 	var channelList agent.ChannelList

// 	for _, ch := range w.Channels {
// 		channelList = append(channelList, &agent.ChannelListing{ID: ch.ID, Type: ch.Type, Name: ch.Name, Topic: ch.Topic})
// 	}

// 	//probably useless
// 	sort.Sort(channelList)

// 	msg := agent.Message{T: agent.MessageChannels}

// 	var err error
// 	msg.M, err = json.Marshal(channelList)
// 	if err != nil {
// 		log.Fatal("failed to marshal ch list : ", err)
// 	}

// 	c.Out <- &msg
// }

// func (w *WebInterface) updateUserList() {
// 	w.ClientsLock.Lock()

// 	msg := &agent.Message{T: agent.MessageUsers}

// 	var userList agent.UserList
// 	for _, wc := range w.Clients {
// 		c := 0
// 		if wc.Channel != nil {
// 			c = wc.Channel.ID
// 		}

// 		userList = append(userList, &agent.User{ID: wc.ID, N: wc.Name, C: c})
// 	}

// 	sort.Sort(userList)

// 	var err error
// 	msg.M, err = json.Marshal(userList)
// 	if err != nil {
// 		log.Fatal("failed to marshal user list: ", err)
// 	}

// 	for _, wc := range w.Clients {
// 		wc.Out <- msg
// 	}

// 	w.ClientsLock.Unlock()
// }

// func (w *WebInterface) quitChannel(c *agent.Client) {
// 	if c.Channel == nil {
// 		return
// 	}

// 	ch := c.Channel

// 	w.ClientsLock.Lock()
// 	ch.Lock()

// 	for _, wc := range ch.Clients {
// 		if len(wc.AudioOut.Tracks) == 0 && wc.ID != c.ID {
// 			continue
// 		}

// 		wc.Out <- &agent.Message{T: agent.MessageQuit, N: c.Name, C: ch.ID}
// 	}

// 	delete(ch.Clients, c.ID)
// 	c.Channel = nil

// 	ch.Unlock()
// 	w.ClientsLock.Unlock()

// 	w.updateUserList()
// }

func DecodeBase64(in string, obj interface{}) {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}
	if compress {
		b = unzip(b)
	}
	err = json.Unmarshal(b, obj)
	if err != nil {
		panic(err)
	}
}

func unzip(in []byte) []byte {
	var b bytes.Buffer
	_, err := b.Write(in)
	if err != nil {
		panic(err)
	}
	r, err := gzip.NewReader(&b)
	if err != nil {
		panic(err)
	}
	res, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return res
}

func EncodeBase64(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}
