export function isBehindOverlay(element: Element | null | undefined): boolean {
  return Boolean(element?.closest("[aria-hidden='true'], [inert]"));
}
