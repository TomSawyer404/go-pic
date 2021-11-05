// go-pic.go - version0.1
// @Note: How it works
//  1. It reads paths from argv[];
//  2. Then construct Post URL to gitee;
//  3. Last get the response from server, display it to stdout
//
// @Author:       MrBanana
// @Date:         2021-8-16

package main

import (
    "fmt"
    "os"
    "io/ioutil"
    "encoding/base64"
    "encoding/json"
    "net/http"
    "net/url"
    "time"
    "math/rand"
    "path/filepath"
)

type Response struct {
    Content struct {
        DownloadUrl string `json:"download_url"`
    }
}

const TOKEN = "2ed36400784dbbfe4515279221aa12e3"

func getRandomString(n int) string {
	randBytes := make([]byte, n)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

func main() {
    PIC_BED_URL := "https://gitee.com/api/v5/repos/tomsawyer404/pic-bed/contents/%2F"

    // Send POST request for every argv[i]
    for i := 1; i < len( os.Args ); i += 1 {
        /// 1. Construct request-body according to gitee's rules
        pic_data, err_read := ioutil.ReadFile( os.Args[i] )
        if err_read != nil {
            fmt.Println("ioutil.ReadFile")
            panic(err_read)
        }
        b64_data := base64.StdEncoding.EncodeToString( []byte(pic_data) )

        post_body := make(url.Values)
        post_body["access_token"] = []string{ TOKEN }
        post_body["owner"] = []string{"tomsawyer404"}
        post_body["repo"] = []string{"pic-bed"}
        post_body["path"] = []string{ "/pic/" }
        post_body["content"] = []string{ b64_data }
        post_body["message"] = []string{ time.Now().Format("2006-01-02 15:04:05") }

        /// 2. Do POST method and receive response
        rand.Seed(time.Now().Unix())
        extension := filepath.Ext( os.Args[i] )
        upload_file_name := getRandomString(8) + extension
        post_response, err := http.PostForm(PIC_BED_URL + upload_file_name, post_body)
        if err != nil {
            fmt.Print("http.PostForm:")
            panic(err)
        }
        defer post_response.Body.Close()

        /// 3. Read Response and find the download url
        body, err := ioutil.ReadAll(post_response.Body)
        if err != nil {
            fmt.Print("ioutil.ReadAll:")
            panic(err)
        }

        var resp Response
        if err := json.Unmarshal(body, &resp); err!= nil {
            fmt.Println("json.Unmarshal", err)
            panic(err)
        }

        //fmt.Println( string(body) )
        fmt.Println( resp.Content.DownloadUrl )
    }

}
