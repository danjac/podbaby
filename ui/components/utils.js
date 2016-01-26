import moment from 'moment';
import sanitizeHtml from 'sanitize-html';
import MobileDetect from 'mobile-detect';

const sanitizeOptions = {
  allowedTags: ['a', 'code', 'em', 'strong', 'b', 'br', 'span', 'img', 'pre'],
  allowedAttributes: {
    a: ['href'],
    span: ['style'],
    img: ['src', 'height', 'width'],
  },
};

export const sanitize = dirty => {
  return {
    __html: sanitizeHtml(dirty, sanitizeOptions),
  };
};

export const formatPubDate = pubDate => moment(pubDate).format('MMMM Do YYYY');

export const formatListenDate = listenedAt => {
  let format = 'MMM D';
  const m = moment(listenedAt);
  if (!m.isSame(new Date(), 'year')) {
    format += ' YYYY';
  }
  return m.format(format);
};

export const isMobile = () => {
  const ua = window && window.navigator && window.navigator.userAgent;
  const md = ua ? new MobileDetect(ua) : null;
  return md && md.mobile() !== null;
};
