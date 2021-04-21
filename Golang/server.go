package main

import (
    rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/RtcTokenBuilder"
    rtmtokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/RtmTokenBuilder"
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
	Channel_name string `json:"ChannelName"`
    Role uint32 `json:"role"`
}

type rtc_string_token_struct struct{
	Uid_rtc_string string `json:"uid"`
	Channel_name string `json:"ChannelName"`
    Role uint32 `json:"role"`
}

type rtm_token_struct struct{
	Uid_rtm string `json:"uid"`
}

var rtc_token string
var rtm_token string
var is_string_uid bool
var int_uid uint32
var string_uid string
var rtm_uid string
var channel_name string

var role_num uint32
var role rtctokenbuilder.Role

func generateRtcToken(is_string_uid bool, int_uid uint32, string_uid string, channelName string){

	appID := "Your_App_ID"
	appCertificate := "Your_Certificate"
	expireTimeInSeconds := uint32(3600)
	currentTimestamp := uint32(time.Now().UTC().Unix())
	expireTimestamp := currentTimestamp + expireTimeInSeconds

	result, err := rtctokenbuilder.BuildTokenWithUID(appID, appCertificate, channelName, int_uid, role, expireTimestamp)
	if (err != nil || is_string_uid == true) {
		fmt.Println(err)
	} else {
		fmt.Printf("Token with uid: %s\n", result)
		fmt.Printf("uid is ", int_uid )
		fmt.Printf("ChannelName is %s\n", channelName)
	}

	result, err = rtctokenbuilder.BuildTokenWithUserAccount(appID, appCertificate, channelName, string_uid, role, expireTimestamp)
	if (err != nil || is_string_uid == false) {
		fmt.Println(err)
	} else {
		fmt.Printf("Token with userAccount: %s\n", result)
		fmt.Printf("uid is %s\n", string_uid)
		fmt.Printf("ChannelName is %s\n", channelName)

	}
	rtc_token = result

}

func generateRtmToken(rtm_uid string){

	appID := "Your_App_ID"
	appCertificate := "Your_Certificate"
	expireTimeInSeconds := uint32(3600)
	currentTimestamp := uint32(time.Now().UTC().Unix())
	expireTimestamp := currentTimestamp + expireTimeInSeconds

	result, err := rtmtokenbuilder.BuildToken(appID, appCertificate, rtm_uid, rtmtokenbuilder.RoleRtmUser, expireTimestamp)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Rtm Token: %s\n", result)

	rtm_token = result

	}
}


func rtcTokenHandler(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS");
    w.Header().Set("Access-Control-Allow-Headers", "*");

    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

    if r.Method != "POST" && r.Method != "OPTIONS" {
        http.Error(w, "Unsupported method. Please check.", http.StatusNotFound)
        return
    }


    var t_int rtc_int_token_struct
	var unmarshalErr *json.UnmarshalTypeError
    int_decoder := json.NewDecoder(r.Body)
	int_err := int_decoder.Decode(&t_int)
	if (int_err == nil) {

                int_uid = t_int.Uid_rtc_int
                is_string_uid = false
                channel_name = t_int.Channel_name
                role_num = t_int.Role
                switch role {
                case 0:
                    role = rtctokenbuilder.RoleAttendee
                case 1:
                    role = rtctokenbuilder.RolePublisher
                case 2:
                    role = rtctokenbuilder.RoleSubscriber
                case 101:
                    role = rtctokenbuilder.RoleAdmin
                }
	}
       if (int_err != nil) {

           if errors.As(int_err, &unmarshalErr){
                   errorResponse(w, "Bad request. Wrong type provided for field " + unmarshalErr.Value  + unmarshalErr.Field + unmarshalErr.Struct, http.StatusBadRequest)
                } else {
                errorResponse(w, "Bad request.", http.StatusBadRequest)
            }
	    return
       }

	generateRtcToken(is_string_uid, int_uid, string_uid, channel_name)
	errorResponse(w, rtc_token, http.StatusOK)
	log.Println(w, r)
}

func rtcStringTokenHandler(w http.ResponseWriter, r *http.Request){

    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS");
    w.Header().Set("Access-Control-Allow-Headers", "*");
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        if r.Method != "POST" && r.Method != "OPTIONS" {
            http.Error(w, "Unsupported method. Please check.", http.StatusNotFound)
            return
        }

        var t_str rtc_string_token_struct
        var unmarshalErr *json.UnmarshalTypeError
        str_decoder := json.NewDecoder(r.Body)
        string_err := str_decoder.Decode(&t_str)

        if (string_err == nil) {
                string_uid = t_str.Uid_rtc_string
                is_string_uid = true
		        channel_name = t_str.Channel_name
                role_num = t_str.Role
                switch role {
                case 0:
                    role = rtctokenbuilder.RoleAttendee
                case 1:
                    role = rtctokenbuilder.RolePublisher
                case 2:
                    role = rtctokenbuilder.RoleSubscriber
                case 101:
                    role = rtctokenbuilder.RoleAdmin
                }
        }

        if (string_err != nil) {
                is_string_uid = false
                if errors.As(string_err, &unmarshalErr){
                   errorResponse(w, "Bad request. Please check your params. ", http.StatusBadRequest)
                } else {
                errorResponse(w, "Bad request.", http.StatusBadRequest)
            }
            return
        }
        // Set a value for int_uid
        int_uid = 0
        generateRtcToken(is_string_uid, int_uid, string_uid, channel_name)
        errorResponse(w, rtc_token, http.StatusOK)
        log.Println(w, r)

}


func rtmTokenHandler(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-Type", "application/json;charset=UTF-8")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS");
        w.Header().Set("Access-Control-Allow-Headers", "*");

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        if r.Method != "POST" && r.Method != "OPTIONS" {
            http.Error(w, "Unsupported method. Please check.", http.StatusNotFound)
            return
        }


        var t_rtm_str rtm_token_struct
        var unmarshalErr *json.UnmarshalTypeError
        str_decoder := json.NewDecoder(r.Body)
        rtm_err := str_decoder.Decode(&t_rtm_str)

        if (rtm_err == nil) {
            rtm_uid = t_rtm_str.Uid_rtm
        }

        if (rtm_err != nil) {
            is_string_uid = false
            if errors.As(rtm_err, &unmarshalErr){
               errorResponse(w, "Bad request. Please check your params.", http.StatusBadRequest)
            } else {
            errorResponse(w, "Bad request.", http.StatusBadRequest)
        }
        return
    }

        generateRtmToken(rtm_uid)
        log.Println(w, r)
}


func errorResponse(w http.ResponseWriter, message string, httpStatusCode int){
	w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["token"] = message
	resp["code"] = strconv.Itoa(httpStatusCode)
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)

}

func main(){
    // Handling routes
    // RTC token from RTC num uid
    http.HandleFunc("/fetch_rtc_token", rtcTokenHandler)
    // RTC token from RTC string uid
    http.HandleFunc("/fetch_rtc_token_string", rtcStringTokenHandler)
    // RTM token from RTM uid
    http.HandleFunc("/fetch_rtm_token", rtmTokenHandler)

    fmt.Printf("Starting server at port 8082\n")

    if err := http.ListenAndServe(":8082", nil); err != nil {
        log.Fatal(err)
    }
}
