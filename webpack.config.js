var path = require('path');
var webpack = require('webpack');
var ExtractTextPlugin = require('extract-text-webpack-plugin');
var UglifyJsPlugin = webpack.optimize.UglifyJsPlugin;

require('es6-promise').polyfill();

var env = process.env.WEBPACK_ENV;

var entry = ['babel-polyfill', './main.js'];

var plugins = [
  new ExtractTextPlugin('[name].css', {
    allChunks: true
  })
];

switch(process.env.WEBPACK_ENV) {
  case 'dev':
    entry.unshift('webpack-dev-server/client?http://localhost:8080');
    entry.unshift('webpack/hot/only-dev-server');
    plugins.unshift(new webpack.HotModuleReplacementPlugin());
    break;
  case 'prod':
    plugins.push(new UglifyJsPlugin({ minimize: true }));
    break;
}

module.exports = {
  context: path.join(__dirname, 'ui'),
  devtool: 'source-map',
  entry: entry,
  output: {
    path: path.join(__dirname, 'dist'),
    filename: "[name].js",
    publicPath: 'http://localhost:8080/static/'
  },
  plugins: plugins,
  module: {
    loaders: [
      {
        test: /\.css$/,
        loader: ExtractTextPlugin.extract('style-loader', 'css-loader')
      },
      {
        test: /\.(png|woff|woff2|eot|ttf|svg)$/,
        loader: 'url-loader?limit=200000'
      },
      {
        test: /\.js$/,
        loader: 'babel-loader',
        query: {
          presets: ['react', 'es2015', 'stage-0']
        },
        include: path.join(__dirname, 'ui'),
        exclude: path.join(__dirname, 'node_modules')
      }
    ]
  },
  resolve: {
    root: path.join(__dirname),
    extensions: ['', '.js'],
    modulesDirectories: ['./ui', './node_modules']
  }

};
