export function getTitle() {
  // generates a document title based on arguments
  return ['Podbaby'].concat(Array.from(arguments)).join(' | ');
}
