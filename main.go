package main

import (
        "os"
        "io"
        "fmt"
        "log"
        "strings"
        "golang.org/x/net/html" 
)


func textify(w io.Writer, n *html.Node) {
    if n.Data == "img" {
        for _, a := range n.Attr {
            if a.Key == "src" {
                io.WriteString(w, "[IMG " + a.Val + "]")
            }
        }
    }
    if n.Data == "a" {
        for _, a := range n.Attr {
            if a.Key == "href" {
                io.WriteString(w, "[HREF " + a.Val + "]")
                return
            }
        }
    }
    if n.Type == html.TextNode {
        io.WriteString(w, n.Data)
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        textify(w, c)
    }
}


func find_tweet_id_and_time(n *html.Node) (tweet_id_and_time string) {
    if n.Type == html.ElementNode && n.Data == "a" {
        for _, a := range n.Attr {
            if a.Key == "href" {
                tweet_id_and_time = a.Val
            }
            if a.Key == "title" {
                tweet_id_and_time = tweet_id_and_time + ", " + a.Val
                return tweet_id_and_time
            }
        }
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {        
        tweet_id_and_time = find_tweet_id_and_time(c)
        if tweet_id_and_time != "" {
            return tweet_id_and_time
        }
    }
    return ""
}


func find_tweet_text(n *html.Node) (tweet_text string) {
    if n.Type == html.ElementNode && n.Data == "p" {
        for _, a := range n.Attr {
            if a.Key == "class" && strings.Contains(a.Val, "tweet-text") {
                var s strings.Builder
                textify(&s, n)
                tweet_text = s.String()
                return tweet_text
            }
        }
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {        
        tweet_text = find_tweet_text(c)
        if tweet_text != "" {
            return tweet_text
        }
    }
    return ""
}


func f(n *html.Node, out *os.File) {
    if n.Type == html.ElementNode && n.Data == "div" {
        for _, a := range n.Attr {
            if a.Key == "class" && a.Val == "content" {
                tweet_id_and_time := find_tweet_id_and_time(n)
                //fmt.Println("tweet_time: " + tweet_time)
                tweet_text := find_tweet_text(n)
                //fmt.Println("tweet_text: " + tweet_text)
                if tweet_id_and_time != "" && tweet_text != "" {
                    out.WriteString(tweet_id_and_time+","+tweet_text+"\n")                    
                }
                break
            }
        }
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        f(c, out)
    }
}


func main() {
    page, err := os.Open(os.Args[1])
    if err != nil {
        log.Fatal(err)
    }
    doc, err := html.Parse(page)
    if err != nil {
        log.Fatal(err)
    }
    out, err := os.OpenFile("tweets.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    out.WriteString("id, time, text\n")
	if err != nil {
		log.Fatal(err)
	}
	f(doc, out)
    fmt.Println("wrote tweets.csv")
    out.Close()
}
