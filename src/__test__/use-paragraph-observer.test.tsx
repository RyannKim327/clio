import { act, cleanup, render } from "@testing-library/react";
import { useParagraphObserver } from "@/hooks/use-paragraph-observer";
import {
  installIntersectionObserverMock,
  MockIntersectionObserver,
} from "./mock-intersection-observer";

function Harness({ threshold }: { threshold?: number }) {
  const [ref, isIntersecting] = useParagraphObserver<HTMLDivElement>(threshold);
  return <div ref={ref} data-testid="node" data-intersecting={isIntersecting.toString()} />;
}

beforeEach(() => installIntersectionObserverMock());
afterEach(() => cleanup());

test("creates an observer for the element with the given threshold", () => {
  const { getByTestId } = render(<Harness threshold={0.5} />);

  expect(MockIntersectionObserver.instances).toHaveLength(1);
  expect(MockIntersectionObserver.instances[0].options?.threshold).toBe(0.5);
  expect(MockIntersectionObserver.instances[0].targets).toContain(getByTestId("node"));
});

test("reports entry and exit from the viewport", () => {
  const { getByTestId } = render(<Harness threshold={0.8} />);
  const node = getByTestId("node");

  expect(node).toBeInTheDocument();
  expect(node.dataset.intersecting).toBe("false");

  act(() => MockIntersectionObserver.instances[0].trigger(true));
  expect(node.dataset.intersecting).toBe("true");

  act(() => MockIntersectionObserver.instances[0].trigger(false));
  expect(node.dataset.intersecting).toBe("false");
});

test("disconnects the observer on unmount", () => {
  const { unmount } = render(<Harness />);
  const instance = MockIntersectionObserver.instances[0];

  unmount();
  expect(instance.disconnected).toBe(true);
});