export function getTitle () {
  // generates a document title based on arguments
  const args = Array.prototype.slice.call(arguments);
  return ["Podbaby"].concat(args).join(" | ");
}
