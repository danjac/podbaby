import moment from 'moment';
import sanitizeHtml from 'sanitize-html';

const sanitizeOptions = {
  allowedTags: ['a', 'code', 'em', 'strong', 'b', 'br', 'span', 'img'],
  allowedAttributes: {
    'a': ['href'],
    'span': ['style'],
    'img': ['src', 'height', 'width']
  }
};

export const sanitize = dirty => {
  return {
    __html: sanitizeHtml(dirty, sanitizeOptions)
  }
};

export const formatPubDate = pubDate => moment(pubDate).format("MMMM Do YYYY");

export function getTitle () {
  const args = Array.prototype.slice.call(arguments);
  return ["Podbaby"].concat(args).join(" | ");
}
