Right now this is just a random collection of things I want to do, in no particular order

- [ ] Real logger
- [ ] Support all the different params (crop dimensions, clamp values, etc) as flags
- [ ] create go script to generate frame images
- [ ] combine frame generation into main script
- [ ] only save frames and cropped frames to disk if debug is enabled
- [ ] Convert OCR function into a multithreaded worker
- [ ] Support progress bar in normal mode
- [ ] support going from link https://www.twitch.tv/videos/1973829739 -> to fully parsing out
- [ ] generate azanlinks page md file
    - [ ] grab meta data on the stream from the given vod id
- [ ] Add way to submit corrections
    - [ ] Include a way to link to the vod timestamp for a specific parsed link
- [ ] Improve url detection
    - [ ] in particular i/L letters are difficult for the OCR to pick up. We can cluster like URLs together and then try resolving the URL and when you find the correct permutation all the like URLs can be mapped to the correct one
- [ ] Support real time link processing (i.e. given a live stream link https://www.twitch.tv/hasanabi)