import { act, cleanup, render } from "@testing-library/react";
import ChapterReader from "@/components/reader/chapter-reader";
import type { c_chapter } from "@/types/book";
import {
  installIntersectionObserverMock,
  MockIntersectionObserver,
} from "./mock-intersection-observer";
import { installMediaStubs } from "./mock-media";

const chapter1: c_chapter[] = [
  { content: "Opening silence" },
  { content: "Track one", music: "one.mp3" },
  { content: "Silence continues" },
  { content: "Track two", music: "two.mp3" },
  { content: "Silence again", end: true },
];

const chapter2: c_chapter[] = [{ content: "New chapter opens silent" }];

let media: ReturnType<typeof installMediaStubs>;
const onPlayingChange = jest.fn();

const props = {
  canPlay: true,
  volume: 1,
  muted: false,
  onPlayingChange,
};

beforeEach(() => {
  media = installMediaStubs();
  onPlayingChange.mockClear();
  installIntersectionObserverMock();
});

afterEach(() => cleanup());

test("plays, continues, switches and stops following the data contract", () => {
  const { container } = render(
    <ChapterReader chapter="chapter 1" paragraphs={chapter1} {...props} />,
  );
  const audio = container.querySelector("audio") as HTMLAudioElement;
  const observers = MockIntersectionObserver.instances;

  expect(media.play).not.toHaveBeenCalled();

  // Music paragraph enters the viewport -> track one starts.
  act(() => observers[1].trigger(true));
  expect(audio.getAttribute("src")).toBe("one.mp3");
  expect(audio.loop).toBe(true);
  expect(media.play).toHaveBeenCalledTimes(1);
  expect(media.pause).not.toHaveBeenCalled();

  // Plain paragraph enters -> the current track continues untouched.
  act(() => observers[2].trigger(true));
  expect(audio.getAttribute("src")).toBe("one.mp3");
  expect(media.play).toHaveBeenCalledTimes(1);
  expect(media.pause).not.toHaveBeenCalled();

  // New music paragraph enters -> the track switches immediately.
  act(() => observers[3].trigger(true));
  expect(audio.getAttribute("src")).toBe("two.mp3");
  expect(media.play).toHaveBeenCalledTimes(2);
  expect(media.pause).toHaveBeenCalledTimes(1);
  expect(audio.currentTime).toBe(0);

  // End paragraph enters -> audio stops and resets.
  act(() => observers[4].trigger(true));
  expect(media.pause).toHaveBeenCalledTimes(2);
  expect(audio.currentTime).toBe(0);

  // Scrolling back up: the music paragraph leaves the viewport and re-enters it.
  act(() => observers[1].trigger(false));
  act(() => observers[1].trigger(true));
  expect(audio.getAttribute("src")).toBe("one.mp3");
  expect(media.play).toHaveBeenCalledTimes(3);
});

test("survives rapid scrolling between music and end paragraphs", () => {
  const { container } = render(
    <ChapterReader chapter="chapter 1" paragraphs={chapter1} {...props} />,
  );
  const audio = container.querySelector("audio") as HTMLAudioElement;
  const observers = MockIntersectionObserver.instances;

  // Each scroll step is a separate intersection frame.
  act(() => observers[1].trigger(true)); // music one starts
  act(() => observers[4].trigger(true)); // end stops it
  act(() => observers[1].trigger(false)); // leave the music paragraph
  act(() => observers[1].trigger(true)); // re-enter -> restarts music one

  expect(audio.getAttribute("src")).toBe("one.mp3");
  expect(media.play).toHaveBeenCalledTimes(2);
  expect(media.pause).toHaveBeenCalledTimes(1);
});

test("stops audio when the chapter changes", () => {
  const { container, rerender } = render(
    <ChapterReader chapter="chapter 1" paragraphs={chapter1} {...props} />,
  );
  const audio = container.querySelector("audio") as HTMLAudioElement;
  act(() => MockIntersectionObserver.instances[1].trigger(true));
  expect(audio.getAttribute("src")).toBe("one.mp3");
  expect(media.play).toHaveBeenCalledTimes(1);

  rerender(<ChapterReader chapter="chapter 2" paragraphs={chapter2} {...props} />);

  expect(media.pause).toHaveBeenCalled();
  expect(audio.currentTime).toBe(0);
});

test("remains silent in a chapter with no music", () => {
  const { container } = render(
    <ChapterReader chapter="chapter 2" paragraphs={chapter2} {...props} />,
  );
  const audio = container.querySelector("audio") as HTMLAudioElement;

  act(() => MockIntersectionObserver.instances[0].trigger(true));
  expect(audio.getAttribute("src")).toBeNull();
  expect(media.play).not.toHaveBeenCalled();
});