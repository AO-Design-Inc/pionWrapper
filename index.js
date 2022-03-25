const pionjs = require('./build/Release/pionWrapper.node');


module.exports = pionjs;

/*
console.log(pionjs.PionOffer('[{\"urls\":[\"stun:stun.l.google.com:19302\"]}]'))
console.log(pionjs.SetRemoteDescription(`{"type":"offer","sdp":"v=0\\r\\no=- 1044856804377246073 1648235642 IN IP4 0.0.0.0\\r\\ns=-\\r\\nt=0 0\\r\\na=fingerprint:sha-256 21:D1:D5:C1:84:65:CC:FF:03:A1:52:7E:46:2A:72:49:98:C0:B0:CB:9D:57:BA:EC:E4:CF:FC:C1:9C:FD:15:66\\r\\na=group:BUNDLE 0\\r\\nm=video 9 UDP/TLS/RTP/SAVPF 125\\r\\nc=IN IP4 0.0.0.0\\r\\na=setup:actpass\\r\\na=mid:0\\r\\na=ice-ufrag:TqDEDeeBlOnRpJzW\\r\\na=ice-pwd:adAHkTweEgwVTBLwQczRTmbmPvgXthTg\\r\\na=rtcp-mux\\r\\na=rtcp-rsize\\r\\na=rtpmap:125 H264/90000\\r\\na=fmtp:125 level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=42e01f\\r\\na=ssrc:3521027110 cname:2a46351b-90b3-425a-8632-fd7188053f0b\\r\\na=ssrc:3521027110 msid:d30ae415-31cd-4d95-9637-4fa34c4b76c2 3c38beb0-f50b-4fea-b93e-e6243baa2d05\\r\\na=ssrc:3521027110 mslabel:d30ae415-31cd-4d95-9637-4fa34c4b76c2\\r\\na=ssrc:3521027110 label:3c38beb0-f50b-4fea-b93e-e6243baa2d05\\r\\na=msid:f80e4218-bc07-43f8-8a8b-e22353002a2c 3c38beb0-f50b-4fea-b93e-e6243baa2d05\\r\\na=sendonly\\r\\n"}`))
*/
