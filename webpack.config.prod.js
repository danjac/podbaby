var path = require('path');
var webpack = require('webpack');
var ExtractTextPlugin = require('extract-text-webpack-plugin');
var UglifyJsPlugin = webpack.optimize.UglifyJsPlugin;
var DefinePlugin = webpack.DefinePlugin;

require('es6-promise').polyfill();

var entry = ['babel-polyfill', './main.js'];

var plugins = [
  new DefinePlugin({
    'process.env.NODE_ENV': JSON.stringify('prod')
  }),
  new ExtractTextPlugin('[name].css', {
    allChunks: true
  }),
  new UglifyJsPlugin({ minimize: true })
];

module.exports = {
  context: path.join(__dirname, 'ui'),
  entry: entry,
  output: {
    path: path.join(__dirname, 'static'),
    filename: "[name].js"
  },
  plugins: plugins,
  module: {
    loaders: [
      {
        test: /\.css$/,
        loader: ExtractTextPlugin.extract('style-loader', 'css-loader')
      },
      {
        test: /\.json$/,
        loader: 'json-loader'
      },
      {
        test: /\.(png|woff|woff2|eot|ttf|svg)/,
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
    extensions: ['', '.js']
  }

};
