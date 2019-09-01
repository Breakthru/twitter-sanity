# twitter-sanity
Convert tweets from a list of twitter users into an easy to read spreadsheet

## Build

```bash
go get golang.org/x/net
go build
```

## Usage

Create a text file called "follows" with a list of the twitter handles you want to download, then run the `./update.sh` bash script.

The script will call `wget` on the raw html page from twitter and call `twitter-sanity` to generate a `.csv` file that can be opened with a spreadsheet program


