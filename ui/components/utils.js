import sanitizeHtml from 'sanitize-html';

const sanitizeOptions = {
  allowedTags: ['a', 'code', 'p', 'em', 'strong', 'b', 'br', 'span'],
  allowedAttributes: {
    'a': ['href'],
    'span': ['style']
  }
};

export const sanitize = dirty => {
  return {
    __html: sanitizeHtml(dirty, sanitizeOptions)
  }
};
