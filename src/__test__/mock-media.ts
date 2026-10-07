export function installMediaStubs(): { play: jest.Mock; pause: jest.Mock; load: jest.Mock; isPaused: () => boolean } {
  let pausedState = true;

  const play = jest.fn(() => {
    pausedState = false;
    return Promise.resolve(undefined);
  });
  const pause = jest.fn(() => {
    pausedState = true;
  });
  const load = jest.fn();

  Object.defineProperty(HTMLMediaElement.prototype, "play", { configurable: true, value: play });
  Object.defineProperty(HTMLMediaElement.prototype, "pause", { configurable: true, value: pause });
  Object.defineProperty(HTMLMediaElement.prototype, "load", { configurable: true, value: load });
  Object.defineProperty(HTMLMediaElement.prototype, "currentTime", {
    configurable: true,
    writable: true,
    value: 0,
  });
  Object.defineProperty(HTMLMediaElement.prototype, "volume", {
    configurable: true,
    writable: true,
    value: 1,
  });
  Object.defineProperty(HTMLMediaElement.prototype, "muted", {
    configurable: true,
    writable: true,
    value: false,
  });
  Object.defineProperty(HTMLMediaElement.prototype, "paused", {
    configurable: true,
    get: () => pausedState,
  });

  return {
    play,
    pause,
    load,
    isPaused: () => pausedState,
  };
}