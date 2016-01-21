var path = require('path');
var webpack = require('webpack');
var HotModuleReplacementPlugin = webpack.HotModuleReplacementPlugin;
var DefinePlugin = webpack.DefinePlugin;
var ExtractTextPlugin = require('extract-text-webpack-plugin');


var entry = [
  'webpack-dev-server/client?http://localhost:8080',
  'webpack/hot/only-dev-server',
  'babel-polyfill',
  './main.js',
];

var plugins = [
  new DefinePlugin({
    'process.env': {
      NODE_ENV: JSON.stringify('dev'),
    },
  }),
  new HotModuleReplacementPlugin(),
  new ExtractTextPlugin('[name].css', {
    allChunks: true,
  }),
];


require('es6-promise').polyfill();

module.exports = {
  context: path.join(__dirname, 'ui'),
  devtool: 'eval',
  entry: entry,
  output: {
    path: path.join(__dirname, 'static'),
    filename: '[name].js',
    publicPath: 'http://localhost:8080/static/',
  },
  plugins: plugins,
  module: {
    loaders: [
      {
        test: /\.css$/,
        loader: ExtractTextPlugin.extract('style-loader', 'css-loader'),
      },
      {
        test: /\.json$/,
        loader: 'json-loader',
      },
      {
        test: /\.(png|woff|woff2|eot|ttf|svg)/,
        loader: 'url-loader?limit=200000',
      },
      {
        test: /\.js$/,
        loaders: ['react-hot', 'babel-loader?presets[]=react,presets[]=es2015,presets[]=stage-0'],
        include: path.join(__dirname, 'ui'),
        exclude: path.join(__dirname, 'node_modules'),
      },
    ],
  },
  resolve: {
    root: path.join(__dirname),
    extensions: ['', '.js'],
  },

};
