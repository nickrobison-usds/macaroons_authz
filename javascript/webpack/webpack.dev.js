const merge = require('webpack-merge');
const path = require('path');
const common = require('./webpack.common.js');

const outputPath = path.join(__dirname, '..', 'dist');

module.exports = merge(common, {
    mode: 'development',
    devtool: 'inline-source-map',
    devServer: {
        contentBase: outputPath,
        port: 30002
    }
})
