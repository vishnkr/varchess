/** Tiny WebAudio beeps for move / capture / check / illegal when sound is enabled. */

let ctx: AudioContext | null = null;
let unlocked = false;

function getCtx(): AudioContext | null {
	if (typeof window === 'undefined') return null;
	if (!ctx) {
		const AC =
			window.AudioContext ||
			(window as unknown as { webkitAudioContext: typeof AudioContext }).webkitAudioContext;
		if (!AC) return null;
		ctx = new AC();
	}
	return ctx;
}

/** Call from a user gesture so browsers allow later WS-triggered beeps. */
export function unlockSfx() {
	const audio = getCtx();
	if (!audio) return;
	if (audio.state === 'suspended') {
		void audio.resume().then(() => {
			unlocked = audio.state === 'running';
		});
	} else {
		unlocked = true;
	}
}

export type SfxKind = 'move' | 'capture' | 'check' | 'illegal' | 'offer' | 'over' | 'join' | 'leave';

function tone(
	audio: AudioContext,
	when: number,
	freq: number,
	dur: number,
	type: OscillatorType,
	peak: number
) {
	const osc = audio.createOscillator();
	const gain = audio.createGain();
	osc.connect(gain);
	gain.connect(audio.destination);
	osc.type = type;
	osc.frequency.setValueAtTime(freq, when);
	gain.gain.setValueAtTime(0.0001, when);
	gain.gain.exponentialRampToValueAtTime(peak, when + 0.01);
	gain.gain.exponentialRampToValueAtTime(0.0001, when + dur);
	osc.start(when);
	osc.stop(when + dur + 0.02);
}

export function playSfx(kind: SfxKind) {
	const audio = getCtx();
	if (!audio) return;

	const start = () => {
		const now = audio.currentTime;
		switch (kind) {
			case 'move':
				tone(audio, now, 520, 0.07, 'sine', 0.2);
				break;
			case 'capture':
				// Low thump + short higher click
				tone(audio, now, 160, 0.11, 'triangle', 0.28);
				tone(audio, now + 0.03, 420, 0.06, 'square', 0.12);
				break;
			case 'check':
				// Sharp ascending double ping
				tone(audio, now, 760, 0.08, 'sine', 0.22);
				tone(audio, now + 0.09, 980, 0.12, 'sine', 0.2);
				break;
			case 'illegal':
				// Dissonant buzz
				tone(audio, now, 220, 0.1, 'sawtooth', 0.14);
				tone(audio, now + 0.04, 185, 0.12, 'sawtooth', 0.12);
				break;
			case 'offer':
				tone(audio, now, 660, 0.14, 'sine', 0.18);
				break;
			case 'over':
				tone(audio, now, 440, 0.28, 'sine', 0.2);
				break;
			case 'join':
				tone(audio, now, 700, 0.09, 'sine', 0.16);
				break;
			case 'leave':
				tone(audio, now, 280, 0.11, 'sine', 0.16);
				break;
		}
	};

	if (audio.state === 'suspended') {
		void audio.resume().then(() => {
			unlocked = audio.state === 'running';
			if (unlocked) start();
		});
		return;
	}
	unlocked = true;
	start();
}
