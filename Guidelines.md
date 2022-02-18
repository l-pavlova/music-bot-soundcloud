## TO RUN:
- export BOT_TOKEN="YOUR BOT TOKEN"
- go run *.go -t $BOT_TOKEN
- soundcloud client id: "YOUR Soundcloud client id"
- soundcloud auth_token: "YOUR Soundcloud auth_token" //add in all headers request


## TO RUN TESTS:
- go test ./... -coverprofile cover.out

## URLS FOR THE API TAKEN STRAIGHT FROM NETWORK TAB (IM SORRY)
- search a song
 https://api-v2.soundcloud.com/search?q=hello&sc_a_id=f55635101bfdf1e8418a36ef0ee8e86f23d9f257&variant_ids=2451&facet=model&user_id=565035-794848-92508-940751&client_id=BmI0Zgypr3dPccFBK9QLjkCpCgvowlzQ&limit=20&offset=0&linked_partitioning=1&app_version=1643966166&app_locale=en

 -stream a song
 https://api-v2.soundcloud.com/media/soundcloud:tracks:59961388/23a40a05-db04-40f1-b58d-dc7a3aaf99ce/stream/hls?client_id=BmI0Zgypr3dPccFBK9QLjkCpCgvowlzQ&track_authorization=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJnZW8iOiJCRyIsInN1YiI6IjEwODIwMjUwNDkiLCJyaWQiOiI0Yzc1YzY2Mi01MWI4LTQ1MDUtOWYxMi0zY2QzZWIxZWEyOGUiLCJpYXQiOjE2NDQxODY2NDh9.13bFTVrAFcIjs6XxD6w1fYUx45J28bqaU1roHuKcxqE