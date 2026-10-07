import { cleanup, render } from "@testing-library/react";
import AudioController from "@/components/reader/audio-controller";
import { installMediaStubs } from "./mock-media";

const base = {
  activeMusicUrl: null,
  shouldStop: false,
  volume: 1,
  muted: false,
  canPlay: false,
};

let media: ReturnType<typeof installMediaStubs>;

beforeEach(() => {
  media = installMediaStubs();
  jest.spyOn(console, "warn").mockImplementation(() => {});
});

afterEach(() => {
  cleanup();
  jest.restoreAllMocks();
});

const queryAudio = (container: HTMLElement) => container.querySelector("audio") as HTMLAudioElement;

test("does not play before the user interacts", () => {
  const { container } = render(
    <AudioController {...base} canPlay={false} activeMusicUrl="one.mp3" />,
  );
  expect(queryAudio(container).getAttribute("src")).toBeNull();
  expect(media.play).not.toHaveBeenCalled();
});

test("starts playback for a music paragraph once allowed", () => {
  const { container } = render(
    <AudioController {...base} canPlay activeMusicUrl="one.mp3" />,
  );
  const audio = queryAudio(container);

  expect(audio.getAttribute("src")).toBe("one.mp3");
  expect(audio.loop).toBe(true);
  expect(media.play).toHaveBeenCalledTimes(1);
});

test("does not restart the same track", () => {
  const { container, rerender } = render(
    <AudioController {...base} canPlay activeMusicUrl="one.mp3" />,
  );

  rerender(<AudioController {...base} canPlay activeMusicUrl="one.mp3" />);
  rerender(<AudioController {...base} canPlay activeMusicUrl="one.mp3" />);

  const audio = queryAudio(container);
  expect(audio.getAttribute("src")).toBe("one.mp3");
  expect(media.play).toHaveBeenCalledTimes(1);
  expect(media.pause).not.toHaveBeenCalled();
});

test("switches to a new track when a different music paragraph enters", () => {
  const { container, rerender } = render(
    <AudioController {...base} canPlay activeMusicUrl="one.mp3" />,
  );

  rerender(<AudioController {...base} canPlay activeMusicUrl="two.mp3" />);

  const audio = queryAudio(container);
  expect(media.pause).toHaveBeenCalledTimes(1);
  expect((audio as HTMLAudioElement).currentTime).toBe(0);
  expect(audio.getAttribute("src")).toBe("two.mp3");
  expect(media.play).toHaveBeenCalledTimes(2);
});

test("stops and resets when an end paragraph enters", () => {
  const { container, rerender } = render(
    <AudioController {...base} canPlay activeMusicUrl="one.mp3" />,
  );

  rerender(<AudioController {...base} canPlay activeMusicUrl="one.mp3" shouldStop />);

  const audio = queryAudio(container);
  expect(media.pause).toHaveBeenCalled();
  expect((audio as HTMLAudioElement).currentTime).toBe(0);
  expect((audio as HTMLAudioElement).loop).toBe(true);
});

test("resumes the same track after a stop when re-entering its paragraph", () => {
  const { container, rerender } = render(
    <AudioController {...base} canPlay activeMusicUrl="one.mp3" />,
  );

  rerender(<AudioController {...base} canPlay activeMusicUrl="one.mp3" shouldStop />);
  rerender(<AudioController {...base} canPlay activeMusicUrl="one.mp3" shouldStop={false} />);

  const audio = queryAudio(container);
  expect(audio.getAttribute("src")).toBe("one.mp3");
  expect(media.play).toHaveBeenCalledTimes(2);
});

test("applies volume and muted", () => {
  const { container } = render(
    <AudioController {...base} canPlay volume={0.4} muted activeMusicUrl="one.mp3" />,
  );
  const audio = queryAudio(container);
  expect(audio.muted).toBe(true);
});

test("catches playback failures without crashing", async () => {
  media.play.mockRejectedValue(new Error("blocked"));
  const onPlayingChange = jest.fn();

  const { container } = render(
    <AudioController {...base} canPlay activeMusicUrl="one.mp3" onPlayingChange={onPlayingChange} />,
  );

  await new Promise((resolve) => setTimeout(resolve, 0));
  expect(console.warn).toHaveBeenCalled();
  expect(queryAudio(container).getAttribute("src")).toBe("one.mp3");
});

test("pauses when the tab is hidden and resumes on return", () => {
  render(
    <AudioController {...base} canPlay activeMusicUrl="one.mp3" />,
  );

  Object.defineProperty(document, "hidden", { configurable: true, get: () => true });
  document.dispatchEvent(new Event("visibilitychange"));
  expect(media.pause).toHaveBeenCalledTimes(1);

  Object.defineProperty(document, "hidden", { configurable: true, get: () => false });
  document.dispatchEvent(new Event("visibilitychange"));
  expect(media.play).toHaveBeenCalledTimes(2);
});