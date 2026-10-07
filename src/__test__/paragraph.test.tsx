import { act, cleanup, render } from "@testing-library/react";
import Paragraph from "@/components/reader/paragraph";
import {
  installIntersectionObserverMock,
  MockIntersectionObserver,
} from "./mock-intersection-observer";

beforeEach(() => installIntersectionObserverMock());
afterEach(() => cleanup());

test("renders paragraph content", () => {
  const { getByText } = render(
    <Paragraph index={0} content="Hello world" isActive={false} onViewportChange={jest.fn()} />,
  );
  expect(getByText("Hello world")).toBeInTheDocument();
});

test("marks music paragraphs and plain paragraphs are not flagged", () => {
  const { getByText } = render(
    <Paragraph index={0} content="Tracked" music="one.mp3" isActive={false} onViewportChange={jest.fn()} />,
  );
  expect(getByText("Audio starts here")).toBeInTheDocument();
});

test("marks end paragraphs", () => {
  const { getByText } = render(
    <Paragraph index={0} content="Done" end isActive={false} onViewportChange={jest.fn()} />,
  );
  expect(getByText("Audio stops here")).toBeInTheDocument();
});

test("applies the active style to the in-viewport paragraph", () => {
  const { container, rerender } = render(
    <Paragraph index={0} content="Active" isActive={false} onViewportChange={jest.fn()} />,
  );
  expect(container.firstElementChild?.className).toContain("border-transparent");

  rerender(<Paragraph index={0} content="Active" isActive onViewportChange={jest.fn()} />);
  expect(container.firstElementChild?.className).toContain("border-border");
});

test("reports viewport entry to the parent with music/end flags", () => {
  const onChange = jest.fn();
  render(
    <Paragraph index={2} content="Music here" music="one.mp3" isActive={false} onViewportChange={onChange} />,
  );

  act(() => MockIntersectionObserver.instances[0].trigger(true));
  expect(onChange).toHaveBeenCalledWith(2, expect.objectContaining({ music: "one.mp3" }));
});

test("does not report after leaving the viewport", () => {
  const onChange = jest.fn();
  render(<Paragraph index={1} content="Plain" isActive={false} onViewportChange={onChange} />);

  act(() => MockIntersectionObserver.instances[0].trigger(true));
  act(() => MockIntersectionObserver.instances[0].trigger(false));
  expect(onChange).toHaveBeenCalledTimes(1);
});