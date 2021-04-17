movies - search yts.mx for torrents

install: `go build && sudo mv movies /usr/local/bin`

```
--query QUERY, -q QUERY            QUERY to search
--rating RATING, -r RATING         minimum imdb user RATING to filter by: 0 to 9 inclusive (default: 0)
--quality QUALITY, --qual QUALITY  file QUALITY to filter by: 720p, 1080p, 2160p, or 3D (default: "1080p")
--genre GENRE, -g GENRE            imdb GENRE from https://www.imdb.com/genre/ to filter by
--sort VALUE, -s VALUE             VALUE to sort by: title, year, rating, peers, seeds, download_count, like_count, or date_added
--order ORDER, -o ORDER            ORDER to order results by: desc or asc
--disable-trackers, --dt           disables trackers in generated magnet links (default: false)
--open                             opens the first search result magnet link (default: false)
--help, -h                         show help (default: false)
```
