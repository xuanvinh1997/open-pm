const SORT_GAP = 65536

export function calculateSortOrder(
  prevSortOrder: number | null,
  nextSortOrder: number | null,
): number {
  if (prevSortOrder !== null && nextSortOrder !== null) {
    return Math.round((prevSortOrder + nextSortOrder) / 2)
  }
  if (prevSortOrder !== null) {
    return prevSortOrder + SORT_GAP
  }
  if (nextSortOrder !== null) {
    return nextSortOrder - SORT_GAP
  }
  return SORT_GAP
}
