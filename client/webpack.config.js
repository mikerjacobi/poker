var webpack = require('webpack');

module.exports = {
  entry: './js/index.jsx',
  output: {
    filename: './js/bundle.js'       
  },
  module: {
    loaders: [
      {
        test: /\.jsx$/,
        loader: 'babel-loader',
        query: {
          presets: ['es2015', 'react']
        }
      }
    ]
  },
  resolve: {
    extensions: ['', '.js', '.jsx'] 
  },
  plugins: [
    new webpack.OldWatchingPlugin()
  ],
};
