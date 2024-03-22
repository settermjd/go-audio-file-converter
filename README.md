# Go Audio File Converter

This is a small web app built with Go (mainly using Go's standard library) that can convert audio files from one audio format to another.
The list of supported codecs is limited to [AAC][aac-url], [MPEG][mpeg-url], [OGG][ogg-url], and [WAV][wav-url], for no other reason than because this is only intended to be a small application, rather than a going concern. 
However, if you'd like to expand on that list, feel free to do so.

## Prerequisites

To use the application, you'll need the following:

- [Go](https://go.dev/dl/). 
  Minimum version of 1.22.0.

To develop the application, you'll need the following, additional, dependencies:

- [npm](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm).
  This is required to install and run the frontend tooling.

## Getting started

To set up the application, clone it locally, change into the cloned project directory, and install the required Go modules, by running the following commands.

```bash
git clone git@github.com:settermjd/go-audio-file-converter.git
cd go-audio-file-converter
go mod tidy
```

[aac-url]: https://en.wikipedia.org/wiki/Advanced_Audio_Coding
[mpeg-url]: https://en.wikipedia.org/wiki/MP3
[ogg-url]: https://en.wikipedia.org/wiki/Ogg
[wav-url]: https://en.wikipedia.org/wiki/WAV