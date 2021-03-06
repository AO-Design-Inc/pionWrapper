#include <node_api.h>
#include <string.h>
#include <stdio.h>
#include "../go-src/pionHandler.h"
#define BUFSIZE 10000

napi_value CSetRemoteDescription(napi_env env, napi_callback_info info) {
	napi_status status;
	size_t argc = 1;
	napi_value argv[1];
	status = napi_get_cb_info(env, info, &argc, argv, NULL, NULL);
	if (status != napi_ok) {
		napi_throw_error(env, NULL, "Failed to parse args");
	}

	char* remoteSDP = malloc(BUFSIZE); 
	size_t result;

	status = napi_get_value_string_utf8(env, argv[0], remoteSDP, BUFSIZE, &result);
	if (status != napi_ok) {
		napi_throw_error(env, NULL, "Invalid string passed in as arg");
	}

	GoSlice remoteSDPGo = {remoteSDP, result, result};

	
	bool remoteSetOk = (bool) SetRemoteDescription(remoteSDPGo);
	napi_value RemoteDescriptionOk;
	status = napi_get_boolean(env, remoteSetOk, &RemoteDescriptionOk); 
	if (status != napi_ok) { 
		napi_throw_error(env, NULL, "bad remote sdp");
	}
	return RemoteDescriptionOk;
}


napi_value StartPionScreenShare(napi_env env, napi_callback_info info) {
	napi_status status;
	size_t argc = 1;
	napi_value argv[1];
	status = napi_get_cb_info(env, info, &argc, argv, NULL, NULL);
	if (status != napi_ok) {
		napi_throw_error(env, NULL, "Failed to parse args");
	}


	char* iceServers = malloc(BUFSIZE); 
	size_t result;

	status = napi_get_value_string_utf8(env, argv[0], iceServers,BUFSIZE, &result);
	if (status != napi_ok) {
		napi_throw_error(env, NULL, "Invalid string passed in as arg");
	}

	GoSlice iceServersGo = {iceServers, result, result};

	
	char* SDPOffer = SpawnConnection(iceServersGo);
	napi_value mySDP;
	status = napi_create_string_utf8(env, SDPOffer, NAPI_AUTO_LENGTH, &mySDP);
	if (status != napi_ok) {
		napi_throw_error(env, NULL, "bad sdpreturned");
	}
	free(SDPOffer);
	free(iceServers);
	return mySDP;
}


napi_value Init(napi_env env, napi_value exports) {
	napi_status status;
	napi_value fn;
	status = napi_create_function(env, NULL , 0, StartPionScreenShare, NULL, &fn);
	if (status != napi_ok) {
		napi_throw_error(env, NULL, "Unable to wrap native function");
	}

	status = napi_set_named_property(env, exports, "PionOffer", fn);
	if (status != napi_ok) {
		napi_throw_error(env, NULL, "Unable to populate exports");
	}

	napi_value fn2;
	status = napi_create_function(env, NULL, 0, CSetRemoteDescription, NULL, &fn2);
	if (status != napi_ok) {
		napi_throw_error(env, NULL, "unable to create function SetRemoteDescription");
	}

	status = napi_set_named_property(env, exports, "SetRemoteDescription", fn2);
	if (status != napi_ok) {
		napi_throw_error(env, NULL, "Unable to create SetRemoteDescription");
	}

	return exports;
}


NAPI_MODULE(NODE_GYP_MODULE_NAME, Init)
