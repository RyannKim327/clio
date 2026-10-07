import { useEffect, useRef } from "react";

interface audioControllerProps {
  activeMusicUrl: string | null
  shouldStop: boolean
  volume: number
  muted: boolean
  canPlay: boolean
  onPlayingChange?: (playing: boolean) => void
}

const FADE_MS = 300;

export default function AudioController({
  activeMusicUrl,
  shouldStop,
  volume,
  muted,
  canPlay,
  onPlayingChange,
}: audioControllerProps) {
  const audioRef = useRef<HTMLAudioElement | null>(null);
  const currentUrlRef = useRef<string | null>(null);
  const volumeRef = useRef(volume);
  const fadeFrameRef = useRef<number | null>(null);
  const onPlayingChangeRef = useRef(onPlayingChange);

  useEffect(() => {
    onPlayingChangeRef.current = onPlayingChange;
  }, [onPlayingChange]);

  useEffect(() => {
    volumeRef.current = volume;
    const audio = audioRef.current;
    if (audio && fadeFrameRef.current === null) {
      audio.volume = volume;
    }
  }, [volume]);

  useEffect(() => {
    const audio = audioRef.current;
    if (audio) {
      audio.muted = muted;
    }
  }, [muted]);

  const cancelFade = () => {
    if (fadeFrameRef.current !== null) {
      cancelAnimationFrame(fadeFrameRef.current);
      fadeFrameRef.current = null;
    }
  };

  const fadeIn = (audio: HTMLAudioElement) => {
    const start = performance.now();
    const step = (now: number) => {
      const t = Math.min((now - start) / FADE_MS, 1);
      audio.volume = Math.round(volumeRef.current * t * 100) / 100;
      if (t < 1) {
        fadeFrameRef.current = requestAnimationFrame(step);
      } else {
        fadeFrameRef.current = null;
        audio.volume = volumeRef.current;
      }
    };
    fadeFrameRef.current = requestAnimationFrame(step);
  };

  const play = (audio: HTMLAudioElement) => {
    audio.volume = 0;
    audio.play().then(() => fadeIn(audio)).catch((e) => {
      audio.volume = volumeRef.current;
      console.warn("Playback failed:", e);
    });
  };

  useEffect(() => {
    const audio = audioRef.current;
    if (!audio || !canPlay) return;

    if (shouldStop) {
      cancelFade();
      audio.pause();
      audio.currentTime = 0;
      audio.volume = volumeRef.current;
      return;
    }

    if (!activeMusicUrl) return;

    const srcMatches = currentUrlRef.current === activeMusicUrl &&
      audio.getAttribute("src") === activeMusicUrl;

    if (!srcMatches) {
      // Music-to-music switch (or first track): pause any previously active
      // track, then load the new one. The very first load has nothing to stop.
      cancelFade();
      if (currentUrlRef.current !== null) {
        audio.pause();
        audio.currentTime = 0;
      }
      audio.src = activeMusicUrl;
      audio.loop = true;
      currentUrlRef.current = activeMusicUrl;
      play(audio);
    } else if (audio.paused) {
      // Same track requested again (e.g. re-entering its paragraph) — resume, do not restart.
      play(audio);
    }
  }, [activeMusicUrl, shouldStop, canPlay]);

  useEffect(() => {
    const audio = audioRef.current;
    if (!audio) return;

    let resumeOnVisible = false;

    const handleVisibility = () => {
      if (document.hidden) {
        if (!audio.paused && canPlay && !shouldStop) {
          resumeOnVisible = true;
          audio.pause();
        }
      } else if (resumeOnVisible && canPlay && !shouldStop && currentUrlRef.current) {
        resumeOnVisible = false;
        audio.play().catch((e) => console.warn("Playback failed:", e));
      }
    };

    document.addEventListener("visibilitychange", handleVisibility);
    return () => document.removeEventListener("visibilitychange", handleVisibility);
  }, [canPlay, shouldStop]);

  useEffect(() => {
    const audio = audioRef.current;
    if (!audio) return;

    const emitPlaying = () => onPlayingChangeRef.current?.(!audio.paused);
    audio.addEventListener("play", emitPlaying);
    audio.addEventListener("pause", emitPlaying);

    return () => {
      audio.removeEventListener("play", emitPlaying);
      audio.removeEventListener("pause", emitPlaying);
    };
  }, []);

  useEffect(() => {
    const audio = audioRef.current;
    return () => {
      if (!audio) return;
      cancelFade();
      audio.pause();
      audio.removeAttribute("src");
      audio.load();
    };
  }, []);

  return <audio ref={audioRef} loop preload="none" className="hidden" aria-hidden />;
}