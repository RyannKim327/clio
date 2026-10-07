import type { c_chapter } from "@/types/book";

export function splitSegmentsByStartEnd(segments: c_chapter[] = []): c_chapter[][] {
  const blocks: c_chapter[][] = [];
  let currentBlock: c_chapter[] = [];
  let inBlock = false;

  for (const item of segments) {
    const isStart = item?.music;
    const isEnd = item?.end;

    if (isStart && !isEnd) {
      // Close any open block first
      if (currentBlock.length > 0) {
        blocks.push(currentBlock);
      }
      currentBlock = [item];
      inBlock = true;
    } else if (isEnd || isStart) {
      currentBlock.push(item);
      blocks.push(currentBlock);
      currentBlock = [];
      inBlock = false;
    } else {
      if (!inBlock && currentBlock.length === 0) {
        // Standalone item before any start
        blocks.push([item]);
      } else {
        currentBlock.push(item);
      }
    }
  }

  // Push any remaining open block
  if (currentBlock.length > 0) {
    blocks.push(currentBlock);
  }

  return blocks;
}
