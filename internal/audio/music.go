package audio

import (
	"sync"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/speaker"
)

// Note represents a musical note with frequency and relative duration
type Note struct {
	Freq     float64 // Frequency in Hz
	Duration float64 // Relative duration (1.0 = base note)
	Silence  bool    // If true, rest instead of playing
}

// MusicPlayer handles background music playback
type MusicPlayer struct {
	mu       sync.Mutex
	enabled  bool
	toggleCh chan bool
	init     sync.Once
}

// Korobeiniki (original Game Boy Tetris theme)
// Based on the classic Russian folk song
var korobeiniki = []Note{
	// Section A (main theme)
	{659.25, 1.0, false}, // E5 - Qua -
	{493.88, 0.5, false}, // B4  - da
	{523.25, 0.5, false}, // C5  - ne
	{587.33, 0.5, false}, // D5  - sta
	{523.25, 0.5, false}, // C5  - ri
	{493.88, 0.5, false}, // B4  - nu
	{440.00, 1.0, false}, // A4  - sha
	{0, 0.5, true},       // rest
	{440.00, 0.5, false}, // A4  - e
	{523.25, 1.0, false}, // C5  - vo
	{659.25, 1.5, false}, // E5  - zyo
	{587.33, 0.5, false}, // D5  - na
	{523.25, 0.5, false}, // C5  - ko
	{493.88, 1.0, false}, // B4  - e
	{523.25, 0.5, false}, // C5  - i
	{440.00, 1.0, false}, // A4  - gni
	{0, 0.5, true},       // rest
	{440.00, 0.5, false}, // A4  - ye
	{329.63, 1.0, false}, // E4  - no
	{392.00, 1.0, false}, // G4  - no
	{440.00, 1.5, false}, // A4  - chi
	{392.00, 0.5, false}, // G4  - ka
	{329.63, 1.0, false}, // E4  - ni

	// Section B
	{659.25, 1.0, false}, // E5  - Ne
	{622.25, 0.5, false}, // D#5 - le
	{587.33, 0.5, false}, // D5  - tyu
	{523.25, 1.0, false}, // C5  - o
	{587.33, 0.5, false}, // D5  - zi
	{493.88, 1.0, false}, // B4  - na
	{440.00, 1.0, false}, // A4  - e
	{415.30, 0.5, false}, // G#4 - vo
	{329.63, 1.0, false}, // E4  - i
	{369.99, 1.0, false}, // F#4 - s'yu
	{440.00, 1.5, false}, // A4  - da
	{392.00, 0.5, false}, // G4  - ni
	{329.63, 1.0, false}, // E4  - cha

	// Section C
	{659.25, 1.0, false}, // E5  - A
	{493.88, 0.5, false}, // B4  - on
	{523.25, 0.5, false}, // C5  - i
	{587.33, 0.5, false}, // D5  - sto
	{698.46, 0.5, false}, // F5  - -e
	{659.25, 0.5, false}, // E5  - go
	{523.25, 1.0, false}, // C5  - ry
	{587.33, 0.5, false}, // D5  - e
	{523.25, 0.5, false}, // C5  - s'yu
	{493.88, 2.0, false}, // B4  - da
	{440.00, 1.0, false}, // A4  - s'yu
	{0, 0.5, true},       // rest
	{440.00, 0.5, false}, // A4  - da
	{523.25, 1.0, false}, // C5  - vo
	{659.25, 1.5, false}, // E5  - zyo
	{587.33, 0.5, false}, // D5  - na
	{523.25, 0.5, false}, // C5  - ko
	{493.88, 1.0, false}, // B4  - e
}

// Base tempo: 400ms per beat (≈150 BPM, authentic Game Boy Tetris speed)
const baseTempo = 400 * time.Millisecond

// ToneStreamer generates a sine wave tone
type ToneStreamer struct {
	sampleRate   beep.SampleRate
	samplePos    int64
	totalSamples int64
	noteFreq     float64
}

// NewToneStreamer creates a new tone streamer for a single note
func NewToneStreamer(sampleRate beep.SampleRate, freq float64, duration time.Duration) *ToneStreamer {
	totalSamples := int64(duration.Seconds() * float64(sampleRate))
	return &ToneStreamer{
		sampleRate:   sampleRate,
		noteFreq:     freq,
		totalSamples: totalSamples,
		samplePos:    0,
	}
}

// Stream generates audio samples
func (t *ToneStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	for i := range samples {
		if t.samplePos >= t.totalSamples {
			return i, false
		}

		// Generate sine wave
		phase := 2 * 3.14159265359 * float64(t.samplePos) * t.noteFreq / float64(t.sampleRate)
		amplitude := 0.3

		// Apply envelope to avoid clicking
		envelope := 1.0
		attack := int64(0.02 * float64(t.totalSamples))
		release := int64(0.08 * float64(t.totalSamples))

		if t.samplePos < attack {
			envelope = float64(t.samplePos) / float64(attack)
		} else if t.samplePos > t.totalSamples-release {
			envelope = float64(t.totalSamples-t.samplePos) / float64(release)
		}

		sample := amplitude * envelope * sin(phase)
		samples[i] = [2]float64{sample, sample}
		t.samplePos++
	}
	return len(samples), true
}

// Err returns any error (none for tone generator)
func (t *ToneStreamer) Err() error {
	return nil
}

// Len returns the total number of samples
func (t *ToneStreamer) Len() int {
	return int(t.totalSamples)
}

// Position returns the current sample position
func (t *ToneStreamer) Position() int {
	return int(t.samplePos)
}

// Seek sets the sample position
func (t *ToneStreamer) Seek(p int) error {
	t.samplePos = int64(p)
	return nil
}

// sin approximates sine function
func sin(x float64) float64 {
	for x > 3.14159265359 {
		x -= 2 * 3.14159265359
	}
	for x < -3.14159265359 {
		x += 2 * 3.14159265359
	}
	return x - x*x*x/6 + x*x*x*x*x/120
}

// NewMusicPlayer creates a new music player
func NewMusicPlayer() *MusicPlayer {
	return &MusicPlayer{
		enabled:  true,
		toggleCh: make(chan bool, 1),
	}
}

// Toggle switches music on/off
func (mp *MusicPlayer) Toggle() {
	mp.mu.Lock()
	mp.enabled = !mp.enabled
	state := mp.enabled
	mp.mu.Unlock()
	select {
	case mp.toggleCh <- state:
	default:
	}
}

// Stop pauses music playback (e.g. on game over). Safe to call multiple times.
func (mp *MusicPlayer) Stop() {
	mp.mu.Lock()
	if !mp.enabled {
		mp.mu.Unlock()
		return
	}
	mp.enabled = false
	mp.mu.Unlock()
	select {
	case mp.toggleCh <- false:
	default:
	}
}

// Start resumes music playback (e.g. on game restart).
func (mp *MusicPlayer) Start() {
	mp.mu.Lock()
	if mp.enabled {
		mp.mu.Unlock()
		return
	}
	mp.enabled = true
	mp.mu.Unlock()
	select {
	case mp.toggleCh <- true:
	default:
	}
}

// IsEnabled returns whether music is currently playing
func (mp *MusicPlayer) IsEnabled() bool {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	return mp.enabled
}

// PlayKorobeiniki plays the Korobeiniki theme with infinite loop.
// The goroutine runs for the lifetime of the process; use Stop/Start to pause/resume.
func (mp *MusicPlayer) PlayKorobeiniki() {
	mp.init.Do(func() {
		speaker.Init(beep.SampleRate(44100), 44100/2) // 500ms buffer
	})

	// Create infinite looping melody using Iterate
	melodyFunc := func() beep.Streamer {
		return createKorobeinikiMelody(baseTempo)
	}
	looped := beep.Iterate(melodyFunc)

	volumeCtrl := &effects.Volume{
		Streamer: looped,
		Base:     2,
		Volume:   0.25,
	}

	ctrl := &beep.Ctrl{Streamer: volumeCtrl, Paused: false}
	speaker.Play(ctrl)

	// Handle toggle signals for the lifetime of the process
	for enabled := range mp.toggleCh {
		speaker.Lock()
		ctrl.Paused = !enabled
		speaker.Unlock()
	}
}

// createKorobeinikiMelody creates the melody sequence with specified tempo
func createKorobeinikiMelody(baseDuration time.Duration) beep.Streamer {
	melody := make([]beep.Streamer, 0, len(korobeiniki))
	for _, note := range korobeiniki {
		if note.Silence {
			silence := NewToneStreamer(beep.SampleRate(44100), 0, time.Duration(float64(baseDuration)*note.Duration))
			melody = append(melody, silence)
		} else {
			streamer := NewToneStreamer(beep.SampleRate(44100), note.Freq, time.Duration(float64(baseDuration)*note.Duration))
			melody = append(melody, streamer)
		}
	}
	return beep.Seq(melody...)
}

// PlayLineClearBeep plays a short beep sound effect for line clears
func (mp *MusicPlayer) PlayLineClearBeep() {
	// Create a pleasant high-pitched beep (880Hz = A5) for 150ms
	beepFreq := 880.0
	beepDuration := 150 * time.Millisecond

	streamer := NewToneStreamer(beep.SampleRate(44100), beepFreq, beepDuration)

	// Apply volume control
	volumeCtrl := &effects.Volume{
		Streamer: streamer,
		Base:     2,
		Volume:   0.3,
	}

	// Play the beep
	speaker.Play(volumeCtrl)

	// Wait for beep to finish
	time.Sleep(beepDuration)
}
