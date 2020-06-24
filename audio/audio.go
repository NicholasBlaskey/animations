package audio

// https://larsimmisch.github.io/pyalsaaudio/terminology.html
// https://github.com/hajimehoshi/oto/blob/master/internal/mux/mux.go
// https://www.codeproject.com/Articles/8295/MPEG-Audio-Frame-Header
// https://github.com/hajimehoshi/go-mp3/blob/master/decode.go
// https://github.com/hajimehoshi/oto/blob/master/player.go
// https://wiki.multimedia.cx/index.php/PCM
//
// Write writes PCM samples to the Player.
//
// The format is as follows:
//   [data]      = [sample 1] [sample 2] [sample 3] ...
//   [sample *]  = [channel 1] ...
//   [channel *] = [byte 1] [byte 2] ...
//
// For example:
//   s1c1b1 s1c1b2 s1c2b1 s1c2b2 s2c1b1 s2c1b2 s2c2b1 s2c2b2
//
// We want to take every two bytes every 4 bytes.
//
// Byte ordering is little endian.
//
// Idea is we need to take divide the sample rate into 60 buckets for fps
// We then create frequency buckets like 20-40hz,... all the way to the max
// Then we put a bargraph per frequency?
//
// Each readCall does a frame
// each frame lasts for
// 26ms (26/1000 of a second). This works out to around 38fps.
// http://www.mp3-converter.com/mp3codec/frames.htm
//
// https://stackoverflow.com/questions/5890499/pcm-audio-amplitude-values
//
// http://www.geosci.usyd.edu.au/users/jboyden/vad/
//
// https://stackoverflow.com/questions/26663494/algorithm-to-draw-waveform-from-audio

func PlayMp3(fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		panic(err)
	}

	// Sample rate, channelNum, bitDepthInBytes, bufferSizeInBytes?
	c, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	p := c.NewPlayer()
	defer p.Close()

	if _, err := io.Copy(p, d); err != nil {
		panic(err)
	}
}
