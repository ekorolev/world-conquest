const path = require('path')
const HtmlWebpackPlugin = require('html-webpack-plugin')

module.exports = {
  mode: 'development',
  entry: './src/index.js',
  output: {
    filename: 'bundle.js',
    path: path.resolve( __dirname, 'dist' )
  },
  module: {
    rules: [
      {
        test: /\.(woff|ttf|otf|eot)$/i,
        loader: 'file-loader',
        options: {
          context: './src/open-iconic-master/fonts'
        }
      },
      { test: /\.(svg)$/i,
        loader: 'file-loader',
        options: {
          name: './src/open-iconic-master/svg/[name].[ext]'
        }
      },
      {
        test: /\.(jpe?g)$/i, 
        use: 'file-loader'
      },
      { test: /\.vue$/, use: 'vue-loader' },
      { test: /\.(scss)$/,
        use: [
          { loader: 'style-loader' }, 
          { loader: 'css-loader' }, 
          { loader: 'postcss-loader', 
            options: {
              plugins: function () { // post css plugins, can be exported to postcss.config.js
                return [
                  require('precss'),
                  require('autoprefixer')
                ];
              }
            }
          },
          { loader: 'sass-loader' }
        ]
      },
      { test: /\.(css)$/, use: [ {loader: 'style-loader'}, {loader: 'css-loader'}]}
    ]
  },
  devServer: {
    contentBase: './dist',
    port: 9000
  },
  plugins: [
    new HtmlWebpackPlugin({
      title: 'Алина&Женя',
      template: 'src/index.html'
    })
  ]
}