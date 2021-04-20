package main

import (
    "RtcTokenBuilder"
    "RtmTokenBuilder"
    "fmt"
    "log"
    "net/http"
    "time"
    "encoding/json"
    "errors"
    "strconv"
)

type rtc_int_token_struct struct{
	Uid_rtc_int uint32 `json:"uid"`
	channelName string `json:"ChannelName"`
}

type rtc_string_token_struct struct{
	Uid_rtc_string string `json:"uid"`
	channelName string `json:"ChannelName"`
}

type rtm_token_struct struct{
	Uid_rtm string `json:"uid"`
}

var rtc_token string
var rtm_token string
var whiteboard_token string
var is_string_uid bool
var int_uid uint32
var string_uid string
var channel_name string

func getRtcToken(is_string_uid bool, int_uid uint32, string_uid string, channelName string){

	appID := "Your App ID"
	appCertificate := "Your App Certificate"
	expireTimeInSeconds := uint32(3600)
	currentTimestamp := uint32(time.Now().UTC().Unix())
	expireTimestamp := currentTimestamp + expireTimeInSeconds

	result, err := rtctokenbuilder.BuildTokenWithUID(appID, appCertificate, channelName, int_uid, rtctokenbuilder.RoleAttendee, expireTimestamp)
	if (err != nil || is_string_uid == true) {
		fmt.Println(err)
	} else {
		fmt.Printf("Token with uid: %s\n", result)
		fmt.Printf("uid is ", int_uid )
		fmt.Printf("ChannelName is ", channelName)
	}

	result, err = rtctokenbuilder.BuildTokenWithUserAccount(appID, appCertificate, channelName, string_uid, rtctokenbuilder.RoleAttendee, expireTimestamp)
	if (err != nil || is_string_uid == false) {
		fmt.Println(err)
	} else {
		fmt.Printf("Token with userAccount: %s\n", result)
		fmt.Printf("uid is", string_uid)
		fmt.Printf("ChannelName is ", channelName)

	}

	rtc_token = result

}

func getRtmToken(uid string){ 

	appID := "Your App ID"
	appCertificate := "Your App Certificate"
	user := "test_user_id"
	expireTimeInSeconds := uint32(3600)
	currentTimestamp := uint32(time.Now().UTC().Unix())
	expireTimestamp := currentTimestamp + expireTimeInSeconds

	result, err := rtmtokenbuilder.BuildToken(appID, appCertificate, user, rtmtokenbuilder.RoleRtmUser, expireTimestamp)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Rtm Token: %s\n", result)

	rtm_token = result

	}
}

func rtcTokenHandler(w http.ResponseWriter, r *http.Request){

	if r.Method != "GET" {
	    http.Error(w, "Unsupported method. Please check.", http.StatusNotFound)
	    return
	}

	headerContentTtype := r.Header.Get("Content-type")
	if headerContentTtype != "application/json" {
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}

  var t_int rtc_int_token_struct
	var unmarshalErr *json.UnmarshalTypeError

	int_decoder := json.NewDecoder(r.Body)
	int_err := int_decoder.Decode(&t_int)

	if (int_err == nil) {

                int_uid = t_int.Uid_rtc_int
                is_string_uid = false
                channel_name = t_int.channelName
	}

       if (int_err != nil) {

           if errors.As(int_err, &unmarshalErr){
                   errorResponse(w, "Bad request. Wrong type provided for field " + unmarshalErr.Value  + unmarshalErr.Field + unmarshalErr.Struct, http.StatusBadRequest)
                } else {
                errorResponse(w, "Bad request.", http.StatusBadRequest)
            }

	    return

       }



	getRtcToken(is_string_uid, int_uid, string_uid, channel_name)

	// fmt.Fprintf(w, rtc_token)
	errorResponse(w, rtc_token, http.StatusOK)

	log.Println(w, r)



}

func rtcStringTokenHandler(w http.ResponseWriter, r *http.Request){

        if r.Method != "GET" {
            http.Error(w, "Unsupported method. Please check.", http.StatusNotFound)
            return
        }

        headerContentTtype := r.Header.Get("Content-type")
        if headerContentTtype != "application/json" {
                errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
                return
        }

        var t_str rtc_string_token_struct
        var unmarshalErr *json.UnmarshalTypeError


        str_decoder := json.NewDecoder(r.Body)


        string_err := str_decoder.Decode(&t_str)

        if (string_err == nil) {
                string_uid = t_str.Uid_rtc_string

                is_string_uid = true
		channel_name = t_str.channelName

    }

        if (string_err != nil) {

                
                is_string_uid = false
                  if errors.As(string_err, &unmarshalErr){
                   errorResponse(w, "Bad request. Wrong type provided for field " + unmarshalErr.Value  + unmarshalErr.Field + unmarshalErr.Struct, http.StatusBadRequest)
                } else {
                errorResponse(w, "Bad request.", http.StatusBadRequest)
            }

            return

    }

    int_uid = 0


        getRtcToken(is_string_uid, int_uid, string_uid, channel_name)

        // fmt.Fprintf(w, rtc_token)
        errorResponse(w, rtc_token, http.StatusOK)

        log.Println(w, r)



}





func rtmTokenHandler(w http.ResponseWriter, r *http.Request){


        if r.Method != "GET" {
            http.Error(w, "Unsupported method. Please check.", http.StatusNotFound)
            return
        }
        getRtmToken()
        fmt.Fprintf(w, rtm_token)
        log.Println(w, r)
}

// TODO: Still hardcoded for whiteboard
func whiteTokenHandler(w http.ResponseWriter, r *http.Request){

        if r.Method != "GET" {
            http.Error(w, "Unsupported method. Please check.", http.StatusNotFound)
            return
        }


	whiteboard_token = "This is a test whiteboardtoken"

        fmt.Fprintf(w, whiteboard_token)

        log.Println(w, r)
}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int){

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["token"] = message
	resp["code"] = strconv.Itoa(httpStatusCode)
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)

}

func main(){
    // Handling routes
    http.HandleFunc("/fetch_rtc_token", rtcTokenHandler)
    http.HandleFunc("/fetch_rtc_token_string", rtcStringTokenHandler)
    http.HandleFunc("/fetch_rtm_token", rtmTokenHandler)
    http.HandleFunc("/fetch_whiteboard_token", whiteTokenHandler)

    fmt.Printf("Starting server at port 8081\n")

    if err := http.ListenAndServe(":8081", nil); err != nil {
        log.Fatal(err)
    }
}
