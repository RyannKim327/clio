# Agentic Coding Specification: Novel Audio Reader App

## Role

You are an expert React frontend engineer. Your task is to build a complete, production-ready novel reading application in React.js where authors can attach background music to specific paragraphs, and audio auto-plays when a reader scrolls to a paragraph containing a `music` URL and stops when a paragraph containing `end: true` scrolls into view.

---

## Goal

Build a React application that:

1. Loads novel data from a JSON file containing multiple chapters.
2. Renders chapters as a continuous (or routed) reading experience.
3. Automatically plays background music when a paragraph with a `music` property scrolls into the viewport.
4. Stops the currently playing audio when a paragraph with `end: true` scrolls into the viewport.
5. Switches to a new audio track when a new paragraph with a different `music` URL scrolls into view while another track is already playing.
6. Handles browser autoplay restrictions gracefully.
7. Provides a clean, readable, book-like UI.

---

## Data Format (Source of Truth)

The novel content is stored in a single JSON file. The top-level keys are chapter names. Each chapter is an array of paragraph objects.

### Paragraph Object Schema

```json
{
  "content": "string (required) — The paragraph text.",
  "music": "string (optional) — Direct URL to an audio file (.mp3, etc.). When this paragraph enters the viewport, this audio begins playing.",
  "end": "boolean (optional) — When true and this paragraph enters the viewport, any currently playing audio must STOP."
}
```

### Example Data Structure

```json
{
  "chapter 1": [
    { "content": "Paragraph with no audio..." },
    { "content": "Another paragraph with no audio..." },
    {
      "music": "https://example.com/audio/osaka-rain-lofi.mp3",
      "content": "Paragraph that triggers audio when scrolled into view..."
    },
    { "content": "Paragraph with no audio. The music from paragraph 3 continues playing." },
    { "content": "Still no audio flag. Music keeps playing." },
    {
      "content": "This paragraph marks where the audio stops.",
      "end": true
    }
  ],
  "chapter 2": [
    { "content": "Chapter 2 opens in silence..." },
    {
      "music": "https://example.com/audio/empty-alley-wind-chimes.mp3",
      "content": "New music starts here..."
    },
    { "content": "Music continues while the reader reads this..." },
    {
      "music": "https://example.com/audio/audio2.mp3",
      "content": "While the previous track is still playing, this paragraph enters view and SWITCHES the audio to a new track."
    },
    { "content": "The new track continues playing here..." },
    {
      "content": "Music stops here. End of chapter.",
      "end": true
    }
  ]
}
```

### Rules Derived from Data

- A paragraph MAY have a `music` URL. If present, that audio BEGINS PLAYING when the paragraph enters the viewport.
- A paragraph WITHOUT a `music` key does NOT change the audio. The previously playing audio continues uninterrupted.
- A paragraph with `end: true` STOPS any currently playing audio when it enters the viewport. It does NOT start new audio, even if `music` is also present (though `music` and `end` should not appear together by design).
- If a paragraph with `music` enters the viewport while another track is already playing, the current track must STOP and the NEW track must begin playing.
- Audio persists across non-music paragraphs until explicitly stopped by an `end: true` paragraph or replaced by a new `music` paragraph.

---

## Audio Behavior (Critical Logic)

### State Machine

The audio system operates in three states:

1. **Idle**: No audio is playing.
2. **Playing**: An audio track is actively playing.
3. **Stopped**: Audio was playing but has been stopped by an `end` paragraph.

### Transitions

| Event | Current State | Action | New State |
|-------|--------------|--------|-----------|
| Paragraph with `music` enters viewport | Idle | Load URL, play audio | Playing |
| Paragraph with `music` enters viewport | Playing | Stop current, load new URL, play | Playing |
| Paragraph with `music` enters viewport | Stopped | Load URL, play audio | Playing |
| Paragraph without `music` or `end` enters viewport | Any | No audio change | Unchanged |
| Paragraph with `end: true` enters viewport | Playing | Pause/stop audio, reset to beginning | Stopped |
| Paragraph with `end: true` enters viewport | Idle | No action | Idle |
| Chapter changes | Any | Stop audio immediately | Idle |

### Music-to-Music Transition Rule

When a new paragraph with `music` enters the viewport and a different track is already playing:

1. Immediately pause the currently playing audio.
2. Reset its `currentTime` to 0.
3. Set the new `src` to the new `music` URL.
4. Call `play()` on the audio element.
5. The audio element should have `loop = true`.

If the new `music` URL is IDENTICAL to the currently playing URL, do NOT restart the track. Let it continue playing seamlessly.

### Important Notes

- Audio should LOOP while in the Playing state.
- Only ONE audio element exists globally. It is reused for all tracks.
- When an `end` paragraph is reached, audio pauses and `currentTime` resets to 0.

---

## Tech Stack

- React 18+ with Functional Components and Hooks
- Vite (preferred) or Create React App
- CSS Modules OR Tailwind CSS for styling
- HTML5 `<audio>` API via React refs (sufficient for MVP)
- Intersection Observer API (native) for scroll detection
- React Router DOM (optional, for multi-chapter routing)

---

## Feature Requirements

### 1. Paragraph Rendering

- Render each paragraph from the chapter array in order inside a readable container.
- Typography should feel like a novel: serif font, generous line height, comfortable margins.
- Paragraphs should have enough vertical spacing to feel distinct.

### 2. Scroll-Triggered Audio

- Each paragraph must be observed using the Intersection Observer API.
- When a paragraph with `music` enters the viewport (threshold: ~0.5), its associated audio must begin playing.
- When a paragraph with `end: true` enters the viewport, any currently playing audio must STOP.
- When a paragraph with `music` enters the viewport while another track is playing, the audio must SWITCH to the new track.
- Non-music, non-end paragraphs must NOT affect audio state.

### 3. Audio Persistence Between Paragraphs

- Audio that starts at a `music` paragraph continues playing as the reader scrolls through subsequent non-music paragraphs.
- Audio only changes when:
  - A new `music` paragraph is reached (switch tracks).
  - An `end` paragraph is reached (stop audio).
  - The chapter changes (stop audio).

### 4. Audio Transitions

- When switching between two different music URLs, the current audio should pause and the new source should load before playing.
- Optional but recommended: Add a brief fade-in / fade-out (200ms–500ms) using volume interpolation to avoid jarring cuts.

### 5. Browser Autoplay Policy Handling

- Browsers block autoplay until the user interacts with the page.
- Implement a full-screen or centered overlay with a "Start Reading" button.
- On click, set a `hasInteracted` state to `true` and hide the overlay.
- Only call `audio.play()` after this interaction has occurred.
- If playback fails, catch the error gracefully and do not crash the app.

### 6. Chapter Navigation

- Provide a Table of Contents (sidebar, drawer, or top navigation) listing all chapters from the JSON keys.
- Clicking a chapter scrolls to its heading or navigates to a route like `/chapter/:chapterId`.
- If using single-page continuous scroll, render all chapters sequentially with sticky chapter headings.
- If using React Router, each chapter should be a route that loads the correct array from the JSON.

### 7. Visual Indicators

- Highlight the paragraph currently in the viewport (e.g., subtle left border or background tint).
- Show a small floating audio indicator (e.g., a sound wave or musical note icon) when audio is actively playing.
- Mark paragraphs that have associated `music` with a subtle icon so the reader knows audio is starting.
- Mark paragraphs with `end: true` with a subtle "end" indicator so the reader knows audio will stop.

---

## Component Architecture

Build the following components. Keep them modular and well-named.

### `App`
- Loads `novel.json`.
- Manages global state: `hasInteracted`, `currentChapter`, `volume`.
- Renders the `StartOverlay` if needed.
- Renders navigation and the main reader area.

### `TableOfContents`
- Props: `chapters` (array of chapter names), `activeChapter`, `onSelect`.
- Renders clickable chapter list.

### `ChapterReader`
- Props: `chapterId` (string), `paragraphs` (array).
- Maps `paragraphs` to `Paragraph` components.
- Manages which paragraph is currently "active" in the viewport.
- Hosts the global `AudioController` for this chapter.
- When a chapter unmounts or changes, ensures audio stops.

### `Paragraph`
- Props: `content`, `music`, `end`, `index`, `isActive`, `onViewportChange`.
- Wraps text in an observable `<p>` or `<div>`.
- Uses a custom hook to report viewport entry/exit to the parent.
- Applies active styling if `isActive` is true.
- Renders a small music icon if `music` is present.
- Renders a small end icon if `end` is true.

### `AudioController`
- Props: `activeMusicUrl` (string or null), `shouldStop` (boolean), `volume` (number), `canPlay` (boolean).
- Holds a single persistent `<audio>` ref.
- Handles source swapping, play, pause, and loop.
- When `shouldStop` is true, pauses audio and resets `currentTime`.
- When `activeMusicUrl` changes to a new URL while audio is playing, stops current track and starts the new one.
- Catches and logs playback errors.
- Optional: implement volume fade on change.

### `StartOverlay`
- Props: `onStart` (function).
- Centered overlay with title, brief description, and a "Start Reading" button.
- Disappears after interaction.

### `useParagraphObserver(threshold)`
- Custom hook that returns `[ref, isIntersecting]`.
- Uses `IntersectionObserver` to watch the referenced element.
- Cleans up observers on unmount.

---

## State Management Logic

### Tracking Active Paragraph

Use a state object in `ChapterReader` to track which paragraph is currently in view:

```js
const [activeParagraph, setActiveParagraph] = useState({
  index: null,
  musicUrl: null,
  end: false,
});
```

When a paragraph enters the viewport:
- If it has `music`, update `activeParagraph` to `{ index, musicUrl: music, end: false }`.
- If it has `end: true`, update `activeParagraph` to `{ index, musicUrl: null, end: true }`.
- If it has neither, and its index matches the current `activeParagraph.index`, clear `activeParagraph` to `{ index: null, musicUrl: null, end: false }`.

Derive two values from `activeParagraph` to pass to `AudioController`:
- `activeMusicUrl`: `activeParagraph.musicUrl` (or `null`).
- `shouldStop`: `activeParagraph.end` (boolean).

### Audio Controller Logic

```js
useEffect(() => {
  if (!canPlay) return;
  const audio = audioRef.current;

  if (shouldStop) {
    audio.pause();
    audio.currentTime = 0;
    return;
  }

  if (activeMusicUrl) {
    if (audio.src !== activeMusicUrl) {
      audio.pause();
      audio.currentTime = 0;
      audio.src = activeMusicUrl;
      audio.loop = true;
      audio.play().catch((e) => console.warn('Playback failed:', e));
    } else if (audio.paused) {
      audio.play().catch((e) => console.warn('Playback failed:', e));
    }
  }
}, [activeMusicUrl, shouldStop, canPlay]);
```

---

## UI / Styling Requirements

- **Container**: Max-width ~680px, centered, padding ~2rem.
- **Font**: Serif stack (`Georgia, 'Times New Roman', serif`). Size ~1.125rem (18px). Line height ~1.8.
- **Paragraph spacing**: `margin-bottom: 1.5rem`.
- **Chapter headings**: Large serif font, generous margin, possibly an ornament or divider below.
- **End marker**: Paragraphs with `end: true` should have extra bottom padding or a centered ornament (e.g., `* * *`) to signal a section boundary.
- **Active paragraph**: Subtle left border (`border-left: 3px solid #888`) or slightly different background (`rgba(0,0,0,0.02)`).
- **Music indicator**: Small note icon (SVG or emoji) near paragraphs with `music`.
- **End indicator**: Small stop icon near paragraphs with `end: true`.
- **Floating audio widget**: Fixed bottom-right or top-right pill showing "Now Playing" status with a mute toggle.

---

## File Structure

```
src/
├── components/
│   ├── App.jsx
│   ├── TableOfContents.jsx
│   ├── ChapterReader.jsx
│   ├── Paragraph.jsx
│   ├── AudioController.jsx
│   └── StartOverlay.jsx
├── hooks/
│   └── useParagraphObserver.js
├── data/
│   └── novel.json
├── styles/
│   └── index.css   (or App.css / modules)
├── utils/
│   └── audio.js    (optional fade helpers)
└── main.jsx
```

---

## Implementation Checklist for the Agent

- [ ] Create React project with Vite.
- [ ] Add `novel.json` to `src/data/` with the exact schema provided.
- [ ] Build `useParagraphObserver` hook with Intersection Observer.
- [ ] Build `Paragraph` component that reports viewport changes upward, including `music` and `end` flags.
- [ ] Build `AudioController` component managing a single `<audio>` element with `activeMusicUrl` and `shouldStop` props.
- [ ] Build `ChapterReader` that tracks active paragraph and derives `activeMusicUrl` and `shouldStop` for `AudioController`.
- [ ] Build `TableOfContents` from JSON keys.
- [ ] Build `StartOverlay` to handle browser autoplay restrictions.
- [ ] Wire everything together in `App.jsx`.
- [ ] Style the app for a clean, readable novel experience.
- [ ] Test scrolling behavior: audio starts at `music` paragraphs, continues through plain paragraphs, switches at new `music` paragraphs, stops at `end` paragraphs.
- [ ] Test rapid scrolling: ensure no overlapping audio or errors.
- [ ] Test chapter navigation: ensure audio resets on chapter change.
- [ ] Test music-to-music transition: ensure seamless track switching without overlap.

---

## Edge Cases to Handle

1. **No music in a chapter**: `AudioController` receives `null` for `activeMusicUrl` and `false` for `shouldStop`. Remains silent.
2. **Same music URL in consecutive paragraphs**: Do not restart the audio; let it continue playing seamlessly.
3. **Rapid scroll past a music paragraph to an end paragraph**: Audio may start briefly, then immediately stop. Ensure no errors.
4. **Missing or broken audio URL**: Catch `audio.play()` promise rejection and fail silently (log to console).
5. **Mobile browsers**: Audio may be blocked more aggressively. The `StartOverlay` interaction is mandatory.
6. **Page visibility change**: If the user switches tabs, audio should pause and resume when they return (use `document.visibilitychange`).
7. **Chapter switch**: Always stop current audio and reset `activeParagraph` when changing chapters.
8. **Music paragraph immediately followed by end paragraph**: Audio starts, then immediately stops when the end paragraph enters view.
9. **Music paragraph immediately followed by another music paragraph**: Current track stops, new track starts immediately.

---

## JSON Data Contract

The agent must treat the following as immutable rules:

- `content` is always present and is a non-empty string.
- `music` is optional. When present, it is a direct URL string to an audio file. It triggers audio playback when its paragraph enters the viewport.
- `end` is optional. When `true`, it stops any currently playing audio when its paragraph enters the viewport.
- `music` and `end` should not appear on the same paragraph by design, but if they do, `end` takes precedence (stop audio, do not play).
- Chapter names are the keys of the top-level JSON object (e.g., `"chapter 1"`, `"chapter 2"`).
- The agent should NOT modify the JSON structure. It must consume it as-is.

---

## Example: How Audio Maps to Paragraphs

### Example 1: Music continues until end

- Paragraph 0: no music, no end
- Paragraph 1: no music, no end
- Paragraph 2: music = "track-a.mp3"
- Paragraph 3: no music, no end
- Paragraph 4: no music, no end
- Paragraph 5: end = true

Behavior:
- Scrolling through paragraphs 0–1: silence.
- Paragraph 2 enters viewport: "track-a.mp3" starts playing.
- Paragraph 3 enters viewport: "track-a.mp3" continues playing (no change).
- Paragraph 4 enters viewport: "track-a.mp3" continues playing (no change).
- Paragraph 5 enters viewport: "track-a.mp3" stops.

### Example 2: Music switches to new track

- Paragraph 0: music = "track-a.mp3"
- Paragraph 1: no music, no end
- Paragraph 2: music = "track-b.mp3"
- Paragraph 3: end = true

Behavior:
- Paragraph 0 enters viewport: "track-a.mp3" plays.
- Paragraph 1 enters viewport: "track-a.mp3" continues.
- Paragraph 2 enters viewport: "track-a.mp3" stops, "track-b.mp3" starts.
- Paragraph 3 enters viewport: "track-b.mp3" stops.

### Example 3: Music switches while playing

- Paragraph 0: music = "track-a.mp3"
- Paragraph 1: music = "track-b.mp3"
- Paragraph 2: end = true

Behavior:
- Paragraph 0 enters viewport: "track-a.mp3" plays.
- Paragraph 1 enters viewport: "track-a.mp3" stops immediately, "track-b.mp3" starts.
- Paragraph 2 enters viewport: "track-b.mp3" stops.

---

## Notes for the Agent

- Do NOT use external state management (Redux, Zustand, etc.) unless explicitly asked. Use React `useState` and `useCallback`.
- Do NOT install heavy audio libraries unless necessary. Native HTML5 `<audio>` is sufficient for the MVP.
- Keep accessibility in mind: use semantic HTML (`<article>`, `<section>`, `<nav>`, `<button>`).
- Ensure the app is responsive and readable on mobile devices.
- Comment complex logic clearly so the human developer can maintain it.
