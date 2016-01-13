if (process.env.NODE_ENV === "dev") {
  module.exports = require('./root.dev');
} else {
  module.exports = require('./root.prod');
}


